package istio

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/kubeerrutils"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/wasme/pkg/cmd/cache"
	"github.com/solo-io/wasme/pkg/deploy"
	envoyfilter "github.com/solo-io/wasme/pkg/deploy/filter"
	"github.com/solo-io/wasme/pkg/pull"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"path/filepath"
	"strings"
)

type WorkloadType string

const (
	WorkloadTypeDeployment = "deployment"
	WorkloadTypeDaemonSet  = "daemonset"

	backupAnnotationPrefix = "wasme_backup."
)

type Workload struct {
	// leave name empty to select ALL workloads in the namespace
	Name      string
	Namespace string
	Type      WorkloadType
}

type Provider struct {
	Ctx         context.Context
	KubeClient  kubernetes.Interface
	IstioClient versionedclient.Interface

	// pulls the image descriptor so we can get the
	// name of the file created by the cache
	Puller pull.DescriptorPuller

	// the target workload to deploy the filter
	Workload Workload
}

// the sidecar annotations required on the pod
var requiredSidecarAnnotations = map[string]string{
	"sidecar.istio.io/userVolume":          `[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache"}}]`,
	"sidecar.istio.io/userVolumeMount":     `[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]`,
	"sidecar.istio.io/interceptionMode":    "TPROXY",
	"sidecar.istio.io/includeInboundPorts": "*",
}

func (p *Provider) setAnnotations(template *v1.PodTemplateSpec) {
	if template.Annotations == nil {
		template.Annotations = map[string]string{}
	}
	for k, v := range requiredSidecarAnnotations {
		// create backups of the existing annotations if they exist
		if currentVal, ok := template.Annotations[k]; ok {
			template.Annotations[backupAnnotationPrefix+k] = currentVal
		}
		template.Annotations[k] = v
	}
}

// runs a function on the workload pod template spec
// selects all workloads in a namespace if workload.Name == ""
func (p *Provider) applyToWorkloadTemplate(do func(spec *v1.PodTemplateSpec)) error {
	switch p.Workload.Type {
	case WorkloadTypeDeployment:
		if p.Workload.Name == "" {
			workloads, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).List(metav1.ListOptions{})
			if err != nil {
				return err
			}
			for _, workload := range workloads.Items {
				do(&workload.Spec.Template)

				if _, err = p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).Update(&workload); err != nil {
					return err
				}
			}
		} else {
			workload, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).Get(p.Workload.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			do(&workload.Spec.Template)

			if _, err = p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).Update(workload); err != nil {
				return err
			}
		}
	case WorkloadTypeDaemonSet:
		if p.Workload.Name == "" {
			workloads, err := p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).List(metav1.ListOptions{})
			if err != nil {
				return err
			}
			for _, workload := range workloads.Items {
				do(&workload.Spec.Template)

				if _, err = p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Update(&workload); err != nil {
					return err
				}
			}
		} else {
			workload, err := p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Get(p.Workload.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			do(&workload.Spec.Template)

			if _, err = p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Update(workload); err != nil {
				return err
			}
		}
	default:
		return errors.Errorf("unknown workload type %v, must be %v or %v", p.Workload.Type, WorkloadTypeDeployment, WorkloadTypeDaemonSet)

	}
	return nil

}

func (p *Provider) ApplyFilter(filter *deploy.Filter) error {
	var labels map[string]string
	// update annotations and grab labels
	err := p.applyToWorkloadTemplate(func(spec *v1.PodTemplateSpec) {
		p.setAnnotations(spec)
		labels = spec.Labels
	})
	if err != nil {
		return errors.Wrap(err, "updating workload annotations")
	}

	logger := logrus.WithFields(logrus.Fields{
		"filter":   filter,
		"workload": p.Workload,
	})

	logger.Info("updated workload sidecar annotations")

	istioEnvoyFilter, err := p.makeIstioEnvoyFilter(
		filter,
		labels,
	)
	if err != nil {
		return err
	}

	if p.Workload.Name == "" {
		// select all workloads in the namespace
		istioEnvoyFilter.Spec.WorkloadSelector = nil
	}

	filterLogger := logger.WithFields(logrus.Fields{
		"envoy_filter_resource": istioEnvoyFilter.Name + "." + istioEnvoyFilter.Namespace,
	})

	_, err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Create(istioEnvoyFilter)
	if err != nil {
		if kubeerrutils.IsAlreadyExists(err) {
			// attempt to update if exists
			existing, err := p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Create(istioEnvoyFilter)
			if err != nil {
				return err
			}

			istioEnvoyFilter.ResourceVersion = existing.ResourceVersion

			_, err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Update(istioEnvoyFilter)

			if err != nil {
				return err
			}

			filterLogger.Info("updated Istio EnvoyFilter resource")
		}
		return err
	} else {
		filterLogger.Info("created Istio EnvoyFilter resource")
	}

	return nil
}

// make Istio EnvoyFilter Custom Resource
func (p *Provider) makeIstioEnvoyFilter(filter *deploy.Filter, labels map[string]string) (*v1alpha3.EnvoyFilter, error) {
	descriptor, err := p.Puller.PullCodeDescriptor(p.Ctx, filter.Image)
	if err != nil {
		return nil, err
	}

	// path to the file in the mounted host volume
	// created by the cache
	filename := filepath.Join(
		"/var/local/lib/wasme-cache",
		cache.Digest2filename(descriptor.Digest),
	)

	wasmFilterConfig := envoyfilter.MakeWasmFilter(filter, envoyfilter.MakeLocalDatasource(filename))

	// here we need to use the gogo proto marshal
	patchValue, err := protoutils.MarshalStruct(wasmFilterConfig)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	spec := networkingv1alpha3.EnvoyFilter{
		WorkloadSelector: &networkingv1alpha3.WorkloadSelector{
			Labels: labels,
		},
		ConfigPatches: []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{{
			ApplyTo: networkingv1alpha3.EnvoyFilter_HTTP_FILTER,
			Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
				Context: networkingv1alpha3.EnvoyFilter_SIDECAR_INBOUND,
				ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
					Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
						FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
							Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
								Name: "envoy.http_connection_manager",
								SubFilter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_SubFilterMatch{
									Name: "envoy.router",
								},
							},
						},
					},
				},
			},
			Patch: &networkingv1alpha3.EnvoyFilter_Patch{
				Operation: networkingv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
				Value:     patchValue,
			},
		}},
	}

	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			// in istio's case, filter ID must be a kube-compliant name
			Name:      istioEnvoyFilterName(p.Workload.Name, filter.ID),
			Namespace: p.Workload.Namespace,
		},
		Spec: spec,
	}, nil
}

func istioEnvoyFilterName(workloadName, filterId string) string {
	if workloadName == "" {
		return filterId
	}
	return workloadName + "-" + filterId
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *deploy.Filter) error {
	logger := logrus.WithFields(logrus.Fields{
		"filter":   filter,
		"workload": p.Workload,
	})

	// remove annotations from workload
	err := p.applyToWorkloadTemplate(func(spec *v1.PodTemplateSpec) {
		for k := range requiredSidecarAnnotations {
			delete(spec.Annotations, k)
		}
		for k, v := range spec.Annotations {
			if strings.HasPrefix(backupAnnotationPrefix, k) {
				key := strings.TrimPrefix(k, backupAnnotationPrefix)
				spec.Annotations[key] = v
				delete(spec.Annotations, key)
			}
		}
	})
	if err != nil {
		return errors.Wrap(err, "removing annotations from workload")
	}

	logger.Info("removed sidecar annotations from workload")

	filterName := istioEnvoyFilterName(p.Workload.Name, filter.ID)

	err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Delete(filterName, nil)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"filter": filterName,
	}).Info("deleted Istio EnvoyFilter resource")

	return nil
}
