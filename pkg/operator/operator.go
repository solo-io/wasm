package operator

import (
	"context"

	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/pkg/ezkube"
	"github.com/solo-io/wasme/pkg/deploy/istio"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/pull"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type filterDeploymentHandler struct {
	ctx context.Context

	kubeClient kubernetes.Interface
	client     ezkube.Ensurer

	puller pull.CodePuller
	cache  istio.Cache
}

func (f *filterDeploymentHandler) Create(obj *v1.FilterDeployment) error {
	return f.deploy(obj)
}

func (f *filterDeploymentHandler) Update(_, obj *v1.FilterDeployment) error {
	return f.deploy(obj)
}

func (f *filterDeploymentHandler) Delete(obj *v1.FilterDeployment) error {
	//for _, dep := range
	return f.undeploy(obj)
}

func (f *filterDeploymentHandler) Generic(obj *v1.FilterDeployment) error {
	// should never be called
	panic("not implemented")
}

func (f *filterDeploymentHandler) deploy(obj *v1.FilterDeployment) error {
	status := v1.FilterDeploymentStatus{
		ObservedGeneration: obj.Generation,
		Workloads:          map[string]*v1.WorkloadStatus{},
	}

	err := f.handleFilter(obj, false, func(workloadMeta metav1.ObjectMeta, err error) {
		workloadStatus := &v1.WorkloadStatus{
			State: v1.WorkloadStatus_Succeeded,
		}
		if err != nil {
			workloadStatus = &v1.WorkloadStatus{
				Reason: err.Error(),
				State:  v1.WorkloadStatus_Failed,
			}
		}
		status.Workloads[obj.Name] = workloadStatus
	})

	if err != nil {
		status.Reason = err.Error()
	}

	obj.Status = status

	return f.client.UpdateStatus(f.ctx, obj)
}

func (f *filterDeploymentHandler) undeploy(obj *v1.FilterDeployment) error {
	return f.handleFilter(obj, true, nil)
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

	deployment, err := getDeployment(obj)
	if err != nil {
		return err
	}

	switch dep := deployment.GetDeploymentType().(type) {
	case *v1.DeploymentSpec_Istio:
		workload := istio.Workload{
			Kind:      dep.Istio.GetKind(),
			Name:      dep.Istio.GetName(),
			Namespace: obj.Namespace,
		}

		provider, err := istio.NewProvider(
			f.ctx,
			f.kubeClient,
			f.client,
			f.puller,
			workload,
			f.cache,
			obj,
			onWorkload,
		)
		if err != nil {
			return err
		}

		if remove {
			return provider.RemoveFilter(filter)
		} else {
			return provider.ApplyFilter(filter)
		}
	default:
		return errors.Errorf("internal error: %T not implemented", deployment)
	}
}
