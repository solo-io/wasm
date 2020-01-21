package operator

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/autopilot/pkg/ezkube"
	"github.com/solo-io/autopilot/pkg/ezkube/mocks"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/deploy/istio"
	providermocks "github.com/solo-io/wasme/pkg/deploy/mocks"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"time"
)

var _ = Describe("FilterDeploymentEventHandler", func() {
	var (
		kubeClient       kubernetes.Interface
		handler          *filterDeploymentHandler
		filterDeployment *v1.FilterDeployment
		client           *mockClient
		provider         *mockProvider
	)
	BeforeEach(func() {
		kubeClient = fake.NewSimpleClientset()

		client = &mockClient{Ensurer: &mocks.Ensurer{}}
		provider = &mockProvider{Provider: &providermocks.Provider{}}

		handler = &filterDeploymentHandler{
			ctx:        context.TODO(),
			kubeClient: kubeClient,
			client:     client,
			cache:      istio.Cache{Name: "cache-name", Namespace: "cache-namespace"},
			makeProviderFn: func(obj *v1.FilterDeployment, puller pull.ImagePuller, onWorkload func(workloadMeta metav1.ObjectMeta, err error)) (deploy.Provider, error) {
				provider.onWorkloadFn = onWorkload
				return provider, nil
			},
		}

		// need to set deletion timestamp or fmt.Sprintf() panics
		d := metav1.NewTime(time.Now())

		filterDeployment = &v1.FilterDeployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "FilterDeployment",
				APIVersion: "wasme.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Generation:        1,
				Name:              "myfilter",
				Namespace:         "bookinfo",
				CreationTimestamp: d,
				DeletionTimestamp: &d,
			},
			Spec: v1.FilterDeploymentSpec{
				Filter: &v1.FilterSpec{
					Image:  "webassemblyhub.io/ilackarms/istio-example:1.4.2",
					Config: `{"name":"hello","value":"world"}`,
				},
				Deployment: &v1.DeploymentSpec{
					DeploymentType: &v1.DeploymentSpec_Istio{Istio: &v1.IstioDeploymentSpec{
						Kind: "Deployment",
					}},
				},
			},
		}
	})
	applyTest := func(applyFunc func (obj *v1.FilterDeployment) error) {
		provider.On("ApplyFilter", filterDeployment.Spec.Filter).Return(nil)
		client.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

		// ensure the status gets set for the workload
		provider.workloadMeta = metav1.ObjectMeta{Name: "test-workload"}
		provider.err = nil

		err := applyFunc(filterDeployment)
		Expect(err).NotTo(HaveOccurred())

		updated := client.updatedObjStatus
		Expect(updated).NotTo(BeNil())

		Expect(updated).To(BeAssignableToTypeOf(&v1.FilterDeployment{}))

		updatedFilter := updated.(*v1.FilterDeployment)
		Expect(updatedFilter.Status).To(Equal(v1.FilterDeploymentStatus{
			ObservedGeneration: 1,
			Workloads: map[string]*v1.WorkloadStatus{
				"test-workload": {State: v1.WorkloadStatus_Succeeded},
			},
		}))
	}
	It("handles create event", func() {
		applyTest(handler.Create)
	})
	It("handles update event", func() {
		applyTest(func(obj *v1.FilterDeployment) error {
			return handler.Update(nil, obj)
		})
	})
	It("handles delete event", func() {
		provider.On("RemoveFilter", filterDeployment.Spec.Filter).Return(nil)
		client.On("UpdateStatus", mock.Anything, mock.Anything).Return(nil)

		err := handler.Delete(filterDeployment)
		Expect(err).NotTo(HaveOccurred())

		updated := client.updatedObjStatus
		Expect(updated).NotTo(BeNil())

		Expect(updated).To(BeAssignableToTypeOf(&v1.FilterDeployment{}))

		updatedFilter := updated.(*v1.FilterDeployment)
		Expect(updatedFilter.Status).To(Equal(v1.FilterDeploymentStatus{
			ObservedGeneration: 1,
			Workloads: map[string]*v1.WorkloadStatus{},
		}))
	})
})

type mockProvider struct {
	workloadMeta metav1.ObjectMeta
	err          error
	onWorkloadFn func(workloadMeta metav1.ObjectMeta, err error)
	*providermocks.Provider
}

func (c *mockProvider) ApplyFilter(f *v1.FilterSpec) error {
	c.onWorkloadFn(c.workloadMeta, c.err)
	return c.Provider.ApplyFilter(f)
}

type mockClient struct {
	updatedObjStatus ezkube.Object
	*mocks.Ensurer
}

func (c *mockClient) UpdateStatus(ctx context.Context, obj ezkube.Object) error {
	c.updatedObjStatus = obj
	return c.Ensurer.UpdateStatus(ctx, obj)
}
