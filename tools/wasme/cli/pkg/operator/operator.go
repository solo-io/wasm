package operator

import (
	"context"
	"time"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/deploy"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pkg/errors"
	"github.com/solo-io/skv2/pkg/ezkube"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/deploy/istio"
	v1 "github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1/controller"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	usernameSecretKey = "username"
	passwordSecretKey = "password"
)

type filterDeploymentHandler struct {
	ctx context.Context

	kubeClient kubernetes.Interface
	client     ezkube.Ensurer

	cache        istio.Cache
	cacheTimeout time.Duration

	// custom overrides for testing
	makePullerFn   func(secretNamespace string, opts *v1.ImagePullOptions) (pull.ImagePuller, error)
	makeProviderFn func(obj *v1.FilterDeployment, puller pull.ImagePuller, onWorkload func(workloadMeta metav1.ObjectMeta, err error)) (deploy.Provider, error)
}

func NewFilterDeploymentHandler(ctx context.Context, kubeClient kubernetes.Interface, client ezkube.Ensurer, cache istio.Cache, cacheTimeout time.Duration) controller.FilterDeploymentEventHandler {
	return &filterDeploymentHandler{ctx: ctx, kubeClient: kubeClient, client: client, cache: cache, cacheTimeout: cacheTimeout}
}

func (f *filterDeploymentHandler) CreateFilterDeployment(obj *v1.FilterDeployment) error {
	return f.deploy(obj)
}

func (f *filterDeploymentHandler) UpdateFilterDeployment(_, obj *v1.FilterDeployment) error {
	return f.deploy(obj)
}

func (f *filterDeploymentHandler) DeleteFilterDeployment(obj *v1.FilterDeployment) error {
	//for _, dep := range
	return f.undeploy(obj)
}

func (f *filterDeploymentHandler) GenericFilterDeployment(obj *v1.FilterDeployment) error {
	// should never be called
	panic("not implemented")
}

func (f *filterDeploymentHandler) deploy(obj *v1.FilterDeployment) error {
	// refresh obj
	if err := f.client.Get(f.ctx, obj); err != nil {
		return err
	}

	status := v1.FilterDeploymentStatus{
		ObservedGeneration: obj.Generation,
		Workloads:          map[string]*v1.WorkloadStatus{},
	}

	setWorkloadStatus := func(workloadMeta metav1.ObjectMeta, err error) {
		workloadStatus := &v1.WorkloadStatus{
			State: v1.WorkloadStatus_Succeeded,
		}
		if err != nil {
			workloadStatus = &v1.WorkloadStatus{
				Reason: err.Error(),
				State:  v1.WorkloadStatus_Failed,
			}
		}
		log.Log.V(1).Info("applied filter to workload", "result", workloadStatus)
		status.Workloads[workloadMeta.Name] = workloadStatus
	}

	err := f.handleFilter(obj, false, setWorkloadStatus)

	if err != nil {
		status.Reason = err.Error()
	}

	obj.Status = status

	if err := f.client.UpdateStatus(f.ctx, obj); err != nil {
		log.Log.Error(err, "failed to update status", "filterdeployment", obj.Name)
	}

	return nil
}

func (f *filterDeploymentHandler) undeploy(obj *v1.FilterDeployment) error {
	// refresh obj
	if err := f.client.Get(f.ctx, obj); err != nil {
		return err
	}

	status := v1.FilterDeploymentStatus{
		ObservedGeneration: obj.Generation,
		Workloads:          map[string]*v1.WorkloadStatus{},
	}

	err := f.handleFilter(obj, true, nil)

	if err != nil {
		status.Reason = err.Error()
	}

	obj.Status = status

	if err := f.client.UpdateStatus(f.ctx, obj); err != nil {
		log.Log.Error(err, "failed to update status", "filterdeployment", obj.Name)
	}

	return nil
}

func getFilter(obj *v1.FilterDeployment) (*v1.FilterSpec, error) {
	filter := obj.Spec.GetFilter()
	if filter == nil {
		return nil, errors.Errorf("must provide spec.filter")
	}
	if filter.Id == "" {
		filter.Id = obj.Name + "." + obj.Namespace
	}
	return filter, nil
}

func getDeployment(obj *v1.FilterDeployment) (*v1.DeploymentSpec, error) {
	deployment := obj.Spec.GetDeployment()
	if deployment == nil {
		return nil, errors.Errorf("must provide spec.deployment")
	}
	return deployment, nil
}

func (f *filterDeploymentHandler) handleFilter(obj *v1.FilterDeployment, remove bool, onWorkload func(workloadMeta metav1.ObjectMeta, err error)) error {
	filter, err := getFilter(obj)
	if err != nil {
		return err
	}

	makePuller := f.makePuller
	if f.makePullerFn != nil {
		makePuller = f.makePullerFn
	}
	puller, err := makePuller(obj.Namespace, filter.GetImagePullOptions())
	if err != nil {
		return err
	}

	makeProvider := f.makeProvider
	if f.makeProviderFn != nil {
		makeProvider = f.makeProviderFn
	}
	deployer, err := makeProvider(obj, puller, onWorkload)
	if err != nil {
		return err
	}

	if remove {
		return deployer.RemoveFilter(filter)
	}

	return deployer.ApplyFilter(filter)
}

func (f *filterDeploymentHandler) makeProvider(obj *v1.FilterDeployment, puller pull.ImagePuller, onWorkload func(workloadMeta metav1.ObjectMeta, err error)) (deploy.Provider, error) {
	deployment, err := getDeployment(obj)
	if err != nil {
		return nil, err
	}

	var provider deploy.Provider
	switch dep := deployment.GetDeploymentType().(type) {
	case *v1.DeploymentSpec_Istio:
		workload := istio.Workload{
			Kind:      dep.Istio.GetKind(),
			Labels:    dep.Istio.GetLabels(),
			Namespace: obj.Namespace,
		}

		provider, err = istio.NewProvider(
			f.ctx,
			f.kubeClient,
			f.client,
			puller,
			workload,
			f.cache,
			obj,
			onWorkload,
			dep.Istio.IstioNamespace,
			"",
			f.cacheTimeout,
			false,
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("internal error: %T not implemented", deployment)
	}
	// deployer sets the root_id on the filter if the user hasn't provided one
	return &deploy.Deployer{
		Ctx:      f.ctx,
		Puller:   puller,
		Provider: provider,
	}, nil
}

func (f *filterDeploymentHandler) makePuller(secretNamespace string, opts *v1.ImagePullOptions) (pull.ImagePuller, error) {
	var username, password string

	if secretName := opts.GetPullSecret(); secretName != "" {
		secret := &kubev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: secretNamespace,
			},
		}
		err := f.client.Get(f.ctx, secret)
		if err != nil {
			return nil, errors.Wrap(err, "missing pull secret")
		}

		if secret.Data == nil {
			return nil, errors.Wrap(err, "secret data is empty")
		}

		u, ok := secret.Data[usernameSecretKey]
		if !ok {
			return nil, errors.Wrapf(err, "secret data missing '%v' key", usernameSecretKey)
		}

		username = string(u)

		p, ok := secret.Data[passwordSecretKey]
		if !ok {
			return nil, errors.Wrapf(err, "secret data missing '%v' key", passwordSecretKey)
		}

		password = string(p)
	}

	resolver, _ := resolver.NewResolver(username, password, opts.GetInsecureSkipVerify(), opts.GetPlainHttp())

	return pull.NewPuller(resolver), nil
}
