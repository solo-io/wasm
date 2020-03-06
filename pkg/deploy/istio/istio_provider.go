package istio

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/solo-io/wasme/pkg/abi"

	"github.com/solo-io/autopilot/pkg/ezkube"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/wasme/pkg/cache"
	envoyfilter "github.com/solo-io/wasme/pkg/deploy/filter"
	"github.com/solo-io/wasme/pkg/pull"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
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

// the sidecar annotations required on the pod
func requiredSidecarAnnotations() map[string]string {
	return map[string]string{
		"sidecar.istio.io/userVolume":      `[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache"}}]`,
		"sidecar.istio.io/userVolumeMount": `[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]`,
	}
}

// applies the filter to all selected workloads and updates the image cache configmap
func (p *Provider) ApplyFilter(filter *v1.FilterSpec) error {

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
			return errors.Errorf("image %v not supported by istio version %v", image.Ref(), istioVersion)
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"image": image.Ref(),
		}).Warnf("no ABI Version found for image, skipping ABI version check")
	}

	if err := p.addImageToCacheConfigMap(filter.Image); err != nil {
		return errors.Wrap(err, "adding image to cache")
	}

	err = p.forEachWorkload(func(meta metav1.ObjectMeta, spec *corev1.PodTemplateSpec) error {
		err := p.applyFilterToWorkload(filter, image, meta, spec)
		if p.OnWorkload != nil {
			p.OnWorkload(meta, err)
		}
		return err
	})
	if err != nil {
		return errors.Wrap(err, "applying filter to workload")
	}

	return nil
}

// applies the filter to the target workload: adds annotations and creates the EnvoyFilter CR
func (p *Provider) applyFilterToWorkload(filter *v1.FilterSpec, image pull.Image, meta metav1.ObjectMeta, spec *corev1.PodTemplateSpec) error {
	p.setAnnotations(spec)
	labels := spec.Labels
	workloadName := meta.Name

	logger := logrus.WithFields(logrus.Fields{
		"filter":   filter,
		"workload": workloadName,
	})

	logger.Info("updated workload sidecar annotations")

	istioEnvoyFilter, err := p.makeIstioEnvoyFilter(
		filter,
		image,
		workloadName,
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

// updates the deployed wasme-cache configmap
// if configmap does not exist (cache not deployed), this will error
func (p *Provider) addImageToCacheConfigMap(image string) error {
	cm, err := p.KubeClient.CoreV1().ConfigMaps(p.Cache.Namespace).Get(p.Cache.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"cache": p.Cache,
		"image": image,
	})

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	images := strings.Split(cm.Data[cache.ImagesKey], "\n")

	for _, existingImage := range images {
		if image == existingImage {
			logger.Info("image is already cached")
			// already exists
			return nil
		}
	}

	images = append(images, image)

	cm.Data[cache.ImagesKey] = strings.Trim(strings.Join(images, "\n"), "\n")

	_, err = p.KubeClient.CoreV1().ConfigMaps(p.Cache.Namespace).Update(cm)
	if err != nil {
		return err
	}

	logger.Info("added image to cache config...")

	if err := p.waitForCacheEvents(image); err != nil {
		return errors.Wrapf(err, "waiting for cache to publish event for image")
	}

	if err := p.cleanupCacheEvents(image); err != nil {
		return errors.Wrapf(err, "cleaning up cache events for image")
	}

	return nil

}

// we want to see a cache event for each cache instance, with each ref
// we can mark the events as processed after receiving
func (p *Provider) waitForCacheEvents(image string) error {

	if p.WaitForCacheTimeout == 0 {
		logrus.Infof("skipping cache events wait")
		return nil
	}

	timeout := time.After(p.WaitForCacheTimeout)
	interval := time.Tick(time.Second)

	logrus.Infof("waiting for event with timeout %v", p.WaitForCacheTimeout)

	cacheDaemonset, err := p.KubeClient.AppsV1().DaemonSets(p.Cache.Namespace).Get(p.Cache.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "getting daemonset for cache %v", p.Cache)
	}

	var eventsErr error
	for {
		select {
		case <-timeout:
			return errors.Errorf("timed out after %s (last err: %v)", p.WaitForCacheTimeout, eventsErr)
		case <-interval:
			events, err := cache.GetImageEvents(p.KubeClient, p.Cache.Namespace, image)
			if err != nil {
				return errors.Wrapf(err, "getting events for image %v", image)
			}
			// expect an event for each cache instance
			var successEvents int32

			for _, evt := range events {
				if evt.Reason == cache.Reason_ImageError {
					logrus.Warnf("event %v was in Error state: %+v", evt.Name, evt)
					continue
				}
				successEvents++
			}

			if successEvents != cacheDaemonset.Status.NumberReady {
				eventsErr = errors.Errorf("expected %v image-ready events for image %v, only found %v", cacheDaemonset.Status.NumberReady, image, successEvents)
				logrus.Warnf("event err: %v", eventsErr)
				continue
			}

			logrus.Debugf("ACK all events for image %v", image)
			return nil
		}
	}
}

func (p *Provider) cleanupCacheEvents(image string) error {
	logrus.Infof("cleaning up cache events for image %v", image)
	events, err := cache.GetImageEvents(p.KubeClient, p.Cache.Namespace, image)
	if err != nil {
		return errors.Wrapf(err, "getting events for image %v", image)
	}

	for _, event := range events {
		if err := p.KubeClient.CoreV1().Events(event.Namespace).Delete(event.Name, nil); err != nil {
			return err
		}
	}

	return nil
}

// runs a function on the workload pod template spec
// selects all workloads in a namespace if workload.Name == ""
func (p *Provider) forEachWorkload(do func(meta metav1.ObjectMeta, spec *corev1.PodTemplateSpec) error) error {
	switch strings.ToLower(p.Workload.Kind) {
	case WorkloadTypeDeployment:
		workloads, err := p.KubeClient.AppsV1().Deployments(p.Workload.Namespace).List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(p.Workload.Labels).String(),
		})
		if err != nil {
			return err
		}
		for _, workload := range workloads.Items {
			if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
				return err
			}

			if err = p.Client.Ensure(p.Ctx, nil, &workload); err != nil {
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
		for _, workload := range workloads.Items {
			if err := do(workload.ObjectMeta, &workload.Spec.Template); err != nil {
				return err
			}

			if err = p.Client.Ensure(p.Ctx, nil, &workload); err != nil {
				return err
			}
		}
	default:
		return errors.Errorf("unknown workload type %v, must be %v or %v", p.Workload.Kind, WorkloadTypeDeployment, WorkloadTypeDaemonSet)
	}

	return nil

}

// set sidecar annotations on the workload
func (p *Provider) setAnnotations(template *corev1.PodTemplateSpec) {
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
func (p *Provider) makeIstioEnvoyFilter(filter *v1.FilterSpec, image pull.Image, workloadName string, labels map[string]string) (*v1alpha3.EnvoyFilter, error) {
	descriptor, err := image.Descriptor()
	if err != nil {
		return nil, err
	}

	// path to the file in the mounted host volume
	// created by the cache
	filename := filepath.Join(
		"/var/local/lib/wasme-cache",
		cache.Digest2filename(descriptor.Digest),
	)

	wasmFilterConfig := envoyfilter.MakeIstioWasmFilter(filter,
		envoyfilter.MakeLocalDatasource(filename),
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
	err := p.forEachWorkload(func(meta metav1.ObjectMeta, spec *corev1.PodTemplateSpec) error {
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
