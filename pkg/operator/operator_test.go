package operator

import (
	"context"
	"time"

	"github.com/solo-io/wasme/pkg/consts"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/autopilot/pkg/ezkube"
	mock_ezkube "github.com/solo-io/autopilot/pkg/ezkube/mocks"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/deploy/istio"
	mock_deploy "github.com/solo-io/wasme/pkg/deploy/mocks"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/pull"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("FilterDeploymentEventHandler", func() {
	var (
		kubeClient       kubernetes.Interface
		handler          *filterDeploymentHandler
		filterDeployment *v1.FilterDeployment
		client           *mockClient
		provider         *mockProvider
		mockCtrl         *gomock.Controller
	)
	BeforeEach(func() {
		kubeClient = fake.NewSimpleClientset()

		mockCtrl = gomock.NewController(GinkgoT())

		client = &mockClient{MockEnsurer: mock_ezkube.NewMockEnsurer(mockCtrl)}
		provider = &mockProvider{MockProvider: mock_deploy.NewMockProvider(mockCtrl)}

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
					Image:  consts.HubDomain + "/ilackarms/istio-test:1.4.2-0",
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
	AfterEach(func() {
		mockCtrl.Finish()
	})
	applyTest := func(applyFunc func(obj *v1.FilterDeployment) error) {
		provider.EXPECT().ApplyFilter(filterDeployment.Spec.Filter).Return(nil)
		client.EXPECT().UpdateStatus(gomock.Any(), gomock.Any()).Return(nil)

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
		provider.EXPECT().RemoveFilter(filterDeployment.Spec.Filter).Return(nil)
		client.EXPECT().UpdateStatus(gomock.Any(), gomock.Any()).Return(nil)

		err := handler.Delete(filterDeployment)
		Expect(err).NotTo(HaveOccurred())

		updated := client.updatedObjStatus
		Expect(updated).NotTo(BeNil())

		Expect(updated).To(BeAssignableToTypeOf(&v1.FilterDeployment{}))

		updatedFilter := updated.(*v1.FilterDeployment)
		Expect(updatedFilter.Status).To(Equal(v1.FilterDeploymentStatus{
			ObservedGeneration: 1,
			Workloads:          map[string]*v1.WorkloadStatus{},
		}))
	})
})

type mockProvider struct {
	workloadMeta metav1.ObjectMeta
	err          error
	onWorkloadFn func(workloadMeta metav1.ObjectMeta, err error)
	*mock_deploy.MockProvider
}

func (c *mockProvider) ApplyFilter(f *v1.FilterSpec) error {
	c.onWorkloadFn(c.workloadMeta, c.err)
	return c.MockProvider.ApplyFilter(f)
}

type mockClient struct {
	updatedObjStatus ezkube.Object
	*mock_ezkube.MockEnsurer
}

func (c *mockClient) UpdateStatus(ctx context.Context, obj ezkube.Object) error {
	c.updatedObjStatus = obj
	return c.MockEnsurer.UpdateStatus(ctx, obj)
}
