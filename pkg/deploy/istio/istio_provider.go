package istio

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/kubeerrutils"
	"github.com/solo-io/go-utils/protoutils"
	cachedeployment "github.com/solo-io/wasme/pkg/cache"
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
)

const (
	WorkloadTypeDeployment = "deployment"
	WorkloadTypeDaemonSet  = "daemonset"

	backupAnnotationPrefix = "wasme-backup."
)

// the target workload to deploy the filter to
// can select all workloads in a namespace
type Workload struct {
	// leave name empty to select ALL workloads in the namespace
	Name      string
	Namespace string
	Type      string
}

// reference to the wasme cache
// we need to update the configmap
type Cache struct {
	Name      string
	Namespace string
}

type Provider struct {
	Ctx         context.Context
	KubeClient  kubernetes.Interface
	IstioClient versionedclient.Interface

	// pulls the image descriptor so we can get the
	// name of the file created by the cache
	Puller pull.CodePuller

	// the target workload to deploy the filter
	Workload Workload

	// reference to the wasme cache
	Cache Cache
}

// the sidecar annotations required on the pod
func requiredSidecarAnnotations() map[string]string {
	return map[string]string{
		"sidecar.istio.io/userVolume":      `[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache"}}]`,
		"sidecar.istio.io/userVolumeMount": `[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]`,
	}
}

// applies the filter to all selected workloads and updates the image cache configmap
func (p *Provider) ApplyFilter(filter *deploy.Filter) error {
	if err := p.addImageToCacheConfigMap(filter.Image); err != nil {
		return errors.Wrap(err, "adding image to cache")
	}

	err := p.forEachWorkload(func(meta metav1.ObjectMeta, spec *v1.PodTemplateSpec) error {
		return p.applyFilterToWorkload(filter, meta, spec)
	})
	if err != nil {
		return errors.Wrap(err, "applying filter to workload")
	}

	return nil
}

// applies the filter to the target workload: adds annotations and creates the EnvoyFilter CR
func (p *Provider) applyFilterToWorkload(filter *deploy.Filter, meta metav1.ObjectMeta, spec *v1.PodTemplateSpec) error {
	p.setAnnotations(spec)
	labels := spec.Labels
	ports := collectContainerPorts(spec)
	workloadName := meta.Name

	logger := logrus.WithFields(logrus.Fields{
		"filter":   filter,
		"workload": workloadName,
		"ports":    ports,
	})

	if len(ports) == 0 {
		logger.Info("no ports detected on workload, skipping")
		return nil
	}

	logger.Info("updated workload sidecar annotations")

	istioEnvoyFilter, err := p.makeIstioEnvoyFilter(
		filter,
		workloadName,
		labels,
	)
	if err != nil {
		return err
	}

	filterLogger := logger.WithFields(logrus.Fields{
		"envoy_filter_resource": istioEnvoyFilter.Name + "." + istioEnvoyFilter.Namespace,
	})

	_, err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Create(istioEnvoyFilter)
	if err != nil {
		if kubeerrutils.IsAlreadyExists(err) {

			// attempt to update if exists
			existing, err := p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Get(istioEnvoyFilter.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			istioEnvoyFilter.ResourceVersion = existing.ResourceVersion

			_, err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Update(istioEnvoyFilter)
			if err != nil {
				return err
			}

			filterLogger.Info("updated Istio EnvoyFilter resource")
		} else {
			return err
		}
	} else {
		filterLogger.Info("created Istio EnvoyFilter resource")
	}

	return nil
}

// updates the deployed wasme-cache configmap
// if configmap does not exist (cache not deployed), this will error
func (p *Provider) addImageToCacheConfigMap(image string) error {
	cm, err := p.KubeClient.CoreV1().ConfigMaps(p.Cache.Namespace).Get(p.Cache.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"cache": p.Cache,
	})

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	images := strings.Split(cm.Data[cachedeployment.ImagesKey], "\n")

	for _, existingImage := range images {
		if image == existingImage {
			logger.Info("image is already cached")
			// already exists
			return nil
		}
	}

	images = append(images, image)

	cm.Data[cachedeployment.ImagesKey] = strings.Trim(strings.Join(images, "\n"), "\n")

	_, err = p.KubeClient.CoreV1().ConfigMaps(p.Cache.Namespace).Update(cm)
	if err != nil {
		return err
	}

	logger.Info("added image to cache")

	return nil

}

// runs a function on the workload pod template spec
// selects all workloads in a namespace if workload.Name == ""
func (p *Provider) forEachWorkload(do func(meta metav1.ObjectMeta, spec *v1.PodTemplateSpec) error) error {
	switch p.Workload.Type {
	case WorkloadTypeDeployment:
		if p.Workload.Name == "" {
			workloads, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).List(metav1.ListOptions{})
			if err != nil {
				return err
			}
			for _, workload := range workloads.Items {
				if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
					return err
				}

				if _, err = p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).Update(&workload); err != nil {
					return err
				}
			}
		} else {
			workload, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).Get(p.Workload.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
				return err
			}

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
				if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
					return err
				}

				if _, err = p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Update(&workload); err != nil {
					return err
				}
			}
		} else {
			workload, err := p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Get(p.Workload.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
				return err
			}

			if _, err = p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).Update(workload); err != nil {
				return err
			}
		}
	default:
		return errors.Errorf("unknown workload type %v, must be %v or %v", p.Workload.Type, WorkloadTypeDeployment, WorkloadTypeDaemonSet)
	}

	return nil

}

// set sidecar annotations on the workload
func (p *Provider) setAnnotations(template *v1.PodTemplateSpec) {
	if template.Annotations == nil {
		template.Annotations = map[string]string{}
	}
	for k, v := range requiredSidecarAnnotations() {
		// create backups of the existing annotations if they exist
		if currentVal, ok := template.Annotations[k]; ok {
			template.Annotations[backupAnnotationPrefix+k] = currentVal
		}
		template.Annotations[k] = v
	}
}

// construct Istio EnvoyFilter Custom Resource
func (p *Provider) makeIstioEnvoyFilter(filter *deploy.Filter, workloadName string, labels map[string]string) (*v1alpha3.EnvoyFilter, error) {
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

	wasmFilterConfig := envoyfilter.MakeHackyIstioWasmFilter(filter,
		// use Filename datasource as Istio doesn't yet support
		// AsyncDatasource
		envoyfilter.MakeFilenameDatasource(filename),
	)

	// here we need to use the gogo proto marshal
	patchValue, err := protoutils.MarshalStruct(wasmFilterConfig)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	// helper func to create a matcher for the target port
	// we need a separate match for each port
	makeMatch := func() *networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch {
		return &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
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
		}
	}

	// each config patch only allows one match, so we
	// have to duplicate the config patch for each port we want
	makeConfigPatch := func(match *networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch) *networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch {
		return &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
			ApplyTo: networkingv1alpha3.EnvoyFilter_HTTP_FILTER,
			Match:   match,
			Patch: &networkingv1alpha3.EnvoyFilter_Patch{
				Operation: networkingv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
				Value:     patchValue,
			},
		}
	}

	// create a config patch for each port
	var configPatches []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch
	configPatches = append(configPatches, makeConfigPatch(makeMatch()))

	spec := networkingv1alpha3.EnvoyFilter{
		WorkloadSelector: &networkingv1alpha3.WorkloadSelector{
			Labels: labels,
		},
		ConfigPatches: configPatches,
	}

	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			// in istio's case, filter ID must be a kube-compliant name
			Name:      istioEnvoyFilterName(workloadName, filter.ID),
			Namespace: p.Workload.Namespace,
		},
		Spec: spec,
	}, nil
}

func istioEnvoyFilterName(workloadName, filterId string) string {
	return workloadName + "-" + filterId
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *deploy.Filter) error {
	logger := logrus.WithFields(logrus.Fields{
		"filter": filter.ID,
	})

	logger.WithFields(logrus.Fields{
		"params": p.Workload,
	}).Info("removing filter from one or more workloads...")

	var workloads []string
	// remove annotations from workload
	err := p.forEachWorkload(func(meta metav1.ObjectMeta, spec *v1.PodTemplateSpec) error {
		// collect the name of the workload so we can delete its filter
		workloads = append(workloads, meta.Name)

		logger := logger.WithFields(logrus.Fields{
			"workload": meta.Name,
		})

		for k := range requiredSidecarAnnotations() {
			delete(spec.Annotations, k)
		}
		logger.Info("removing sidecar annotations from workload")

		// restore backup annotations
		for k, v := range spec.Annotations {
			if strings.HasPrefix(backupAnnotationPrefix, k) {
				key := strings.TrimPrefix(k, backupAnnotationPrefix)
				spec.Annotations[key] = v
				delete(spec.Annotations, key)
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "removing annotations from workload")
	}

	for _, workloadName := range workloads {

		filterName := istioEnvoyFilterName(workloadName, filter.ID)

		err = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(p.Workload.Namespace).Delete(filterName, nil)
		if err != nil {
			return err
		}

		logger.WithFields(logrus.Fields{
			"filter": filterName,
		}).Info("deleted Istio EnvoyFilter resource")
	}

	return nil
}

func collectContainerPorts(spec *v1.PodTemplateSpec) []uint32 {
	var ports []uint32
	for _, container := range spec.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, uint32(port.ContainerPort))
		}
	}
	return ports
}
