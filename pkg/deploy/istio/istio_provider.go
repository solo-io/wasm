package istio

import (
	"context"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/solo-io/wasme/pkg/abi"
	"github.com/solo-io/wasme/pkg/consts"

	"github.com/solo-io/skv2/pkg/ezkube"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/protoutils"
	envoyfilter "github.com/solo-io/wasme/pkg/deploy/filter"
	"github.com/solo-io/wasme/pkg/pull"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
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
	Labels    map[string]string
	Namespace string
	Kind      string
}

// reference to the wasme cache
// we need to update the configmap
type Cache struct {
	Name      string
	Namespace string
}

type Provider struct {
	Ctx        context.Context
	KubeClient kubernetes.Interface
	Client     ezkube.Ensurer

	// pulls the image descriptor so we can get the
	// name of the file created by the cache
	Puller pull.ImagePuller

	// the target workload to deploy the filter
	Workload Workload

	// reference to the wasme cache
	Cache Cache

	// set owner references on created Filters with this parent object
	// if it's nil, they will not have an owner reference set
	ParentObject ezkube.Object

	// Callback to the caller when for when the istio provider
	// updates a workload.
	// err != nil in the case that update failed
	OnWorkload func(workloadMeta metav1.ObjectMeta, err error)

	// namespace of the istio control plane
	// Provider will use this to determine the installed version of istio
	// for abi compatibility
	// defaults to istio-system
	IstioNamespace string

	// if non-zero, wait for cache events to be populated with this timeout before
	// creating istio EnvoyFilters.
	// set to zero to skip the check
	WaitForCacheTimeout time.Duration
}

func NewProvider(ctx context.Context, kubeClient kubernetes.Interface, client ezkube.Ensurer, puller pull.ImagePuller, workload Workload, cache Cache, parentObject ezkube.Object, onWorkload func(workloadMeta metav1.ObjectMeta, err error), istioNamespace string, cacheTimeout time.Duration) (*Provider, error) {

	// ensure istio types are added to scheme
	if err := v1alpha3.AddToScheme(client.Manager().GetScheme()); err != nil {
		return nil, err
	}

	return &Provider{
		Ctx:                 ctx,
		KubeClient:          kubeClient,
		Client:              client,
		Puller:              puller,
		Workload:            workload,
		Cache:               cache,
		ParentObject:        parentObject,
		OnWorkload:          onWorkload,
		IstioNamespace:      istioNamespace,
		WaitForCacheTimeout: cacheTimeout,
	}, nil
}

// applies the filter to all selected workloads and updates the image cache configmap
func (p *Provider) ApplyFilter(filter *v1.FilterSpec) error {
	var warning error
	image, err := p.Puller.Pull(p.Ctx, filter.Image)
	if err != nil {
		return err
	}

	cfg, err := image.FetchConfig(p.Ctx)
	if err != nil {
		return err
	}

	abiVersions := cfg.AbiVersions

	if len(abiVersions) > 0 {
		istioVersion, err := p.getIstioVersion()
		if err != nil {
			return err
		}

		if err := abi.DefaultRegistry.ValidateIstioVersion(abiVersions, istioVersion); err != nil {
			warning = errors.Errorf("warning: image %v may not be supported by istio version %v", image.Ref(), istioVersion)
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"image": image.Ref(),
		}).Warnf("no ABI Version found for image, skipping ABI version check")
	}

	err = p.forEachWorkload(func(meta metav1.ObjectMeta, typ metav1.TypeMeta, spec *corev1.PodTemplateSpec) error {
		err := p.applyFilterToWorkload(filter, image, meta, typ, spec)
		if p.OnWorkload != nil {
			p.OnWorkload(meta, err)
		}
		return err
	})
	if err != nil {
		return errors.Wrap(err, "applying filter to workload")
	}

	return warning
}

// applies the filter to the target workload: adds annotations and creates the EnvoyFilter CR
func (p *Provider) applyFilterToWorkload(filter *v1.FilterSpec, image pull.Image, meta metav1.ObjectMeta, typ metav1.TypeMeta, spec *corev1.PodTemplateSpec) error {
	labels := spec.Labels
	workloadName := meta.Name

	logger := logrus.WithFields(logrus.Fields{
		"filter":   filter,
		"workload": workloadName,
	})

	clusterFilter := p.makeClusterFilter()

	err := p.Client.Ensure(p.Ctx, nil, clusterFilter)
	if err != nil {
		return err
	}

	istioEnvoyFilter, err := p.makeIstioEnvoyFilter(
		filter,
		image,
		meta,
		typ,
		labels,
	)
	if err != nil {
		return err
	}

	filterLogger := logger.WithFields(logrus.Fields{
		"envoy_filter_resource": istioEnvoyFilter.Name + "." + istioEnvoyFilter.Namespace,
	})

	err = p.Client.Ensure(p.Ctx, p.ParentObject, istioEnvoyFilter)
	if err != nil {
		return err
	}
	filterLogger.Info("created Istio EnvoyFilter resource")

	return nil
}

// runs a function on the workload pod template spec
// selects all workloads in a namespace if workload.Name == ""
func (p *Provider) forEachWorkload(do func(meta metav1.ObjectMeta, typ metav1.TypeMeta, spec *corev1.PodTemplateSpec) error) error {
	switch strings.ToLower(p.Workload.Kind) {
	case WorkloadTypeDeployment:
		workloads, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(p.Workload.Labels).String(),
		})
		if err != nil {
			return err
		}
		typeMeta := metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		}
		for _, workload := range workloads.Items {
			if err := do(workload.ObjectMeta, typeMeta, &workload.Spec.Template); err != nil {
				return err
			}
		}
	case WorkloadTypeDaemonSet:
		workloads, err := p.KubeClient.AppsV1().DaemonSets(p.Workload.Namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(p.Workload.Labels).String(),
		})
		if err != nil {
			return err
		}
		typeMeta := metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		}
		for _, workload := range workloads.Items {
			if err := do(workload.ObjectMeta, typeMeta, &workload.Spec.Template); err != nil {
				return err
			}
		}
	default:
		return errors.Errorf("unknown workload type %v, must be %v or %v", p.Workload.Kind, WorkloadTypeDeployment, WorkloadTypeDaemonSet)
	}

	return nil

}
func (p *Provider) clusterName() string {
	return "wasme-cache-cluster-" + p.Cache.Name + "-" + p.Cache.Namespace
}

// construct Istio EnvoyFilter Custom Resource
func (p *Provider) makeIstioEnvoyFilter(filter *v1.FilterSpec, image pull.Image, meta metav1.ObjectMeta, typ metav1.TypeMeta, labels map[string]string) (*v1alpha3.EnvoyFilter, error) {
	descriptor, err := image.Descriptor()
	if err != nil {
		return nil, err
	}
	workloadName := meta.Name
	sha := strings.TrimPrefix(string(descriptor.Digest), "sha256:")
	// path to the file in the mounted host volume
	// created by the cache
	clusterName := p.clusterName()
	wasmFilterConfig := envoyfilter.MakeIstioWasmFilter(filter,
		envoyfilter.MakeRemoteDataSource("http://"+clusterName+"/"+image.Ref(), clusterName, sha), // get cluster name nad filter hash
	)

	// here we need to use the gogo proto marshal
	patchValue, err := protoutils.MarshalStruct(wasmFilterConfig)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

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
			Name:      istioEnvoyFilterName(workloadName, filter.Id),
			Namespace: p.Workload.Namespace,
		},
		Spec: spec,
	}, nil
}

// construct Istio EnvoyFilter Custom Resource
func (p *Provider) makeClusterFilter() *v1alpha3.EnvoyFilter {
	// path to the file in the mounted host volume
	// created by the cache
	clusterName := p.clusterName()

	cluster := &envoy_api_v2.Cluster{
		Name:                 clusterName,
		ClusterDiscoveryType: &envoy_api_v2.Cluster_Type{Type: envoy_api_v2.Cluster_STRICT_DNS},
		ConnectTimeout:       &duration.Duration{Seconds: 3},
		Hosts: []*envoy_api_v2_core.Address{
			{
				Address: &envoy_api_v2_core.Address_SocketAddress{
					SocketAddress: &envoy_api_v2_core.SocketAddress{
						Address: p.Cache.Name + "." + p.Cache.Namespace + ".svc.cluster.local", // do we need the suffix svc.cluster.local?
						PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
							PortValue: consts.CachePort,
						},
					},
				},
			},
		},
	}

	clusterValue, err := protoutils.MarshalStruct(cluster)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	makeClusterConfigPatch := func() *networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch {
		return &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
			Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
				Context: networkingv1alpha3.EnvoyFilter_ANY,
				ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
					Cluster: &networkingv1alpha3.EnvoyFilter_ClusterMatch{},
				},
			},
			ApplyTo: networkingv1alpha3.EnvoyFilter_CLUSTER,
			Patch: &networkingv1alpha3.EnvoyFilter_Patch{
				// use merge in case there is more than one.
				Operation: networkingv1alpha3.EnvoyFilter_Patch_ADD,
				Value:     clusterValue,
			},
		}
	}

	var configPatches []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch
	configPatches = append(configPatches, makeClusterConfigPatch())

	spec := networkingv1alpha3.EnvoyFilter{
		ConfigPatches: configPatches,
	}

	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Cache.Name,
			Namespace: p.Workload.Namespace,
		},
		Spec: spec,
	}
}

func istioEnvoyFilterName(workloadName, filterId string) string {
	return workloadName + "-" + filterId
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *v1.FilterSpec) error {
	logger := logrus.WithFields(logrus.Fields{
		"filter": filter.Id,
	})

	logger.WithFields(logrus.Fields{
		"params": p.Workload,
	}).Info("removing filter from one or more workloads...")

	var workloads []string
	// remove annotations from workload
	err := p.forEachWorkload(func(meta metav1.ObjectMeta, typ metav1.TypeMeta, spec *corev1.PodTemplateSpec) error {
		// collect the name of the workload so we can delete its filter
		workloads = append(workloads, meta.Name)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "removing annotations from workload")
	}

	if p.ParentObject != nil {
		// no need to remove the istio filters as they will be garbage collected
		return nil
	}

	for _, workloadName := range workloads {

		filterName := istioEnvoyFilterName(workloadName, filter.Id)

		err = p.Client.Delete(p.Ctx, &v1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: p.Workload.Namespace,
				Name:      filterName,
			},
		})
		if err != nil {
			return err
		}

		logger.WithFields(logrus.Fields{
			"filter": filterName,
		}).Info("deleted Istio EnvoyFilter resource")
	}

	return nil
}

func (p *Provider) getIstioVersion() (string, error) {
	inspector := &versionInspector{
		kube:           p.KubeClient,
		istioNamespace: p.IstioNamespace,
	}
	return inspector.GetIstioVersion()
}
