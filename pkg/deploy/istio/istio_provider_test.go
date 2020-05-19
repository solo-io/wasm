package istio

import (
	"context"
	"fmt"

	"github.com/solo-io/wasme/pkg/consts/test"

	"github.com/solo-io/wasme/pkg/consts"

	"github.com/golang/mock/gomock"
	mock_ezkube "github.com/solo-io/skv2/pkg/ezkube/mocks"

	"github.com/solo-io/wasme/pkg/resolver"

	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/model"
	"github.com/solo-io/wasme/pkg/pull"

	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/randutils"
	"github.com/solo-io/skv2/pkg/ezkube"
	aptest "github.com/solo-io/skv2/test"
	wasmev1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"istio.io/api/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("IstioProvider", func() {
	var (
		kube   kubernetes.Interface
		client ezkube.Ensurer
		ns     string
		cache  Cache
		filter = &wasmev1.FilterSpec{
			Id:     "filter-id",
			Config: `{"filter":"config"}`,
			Image:  "filter/image:v1",
			RootID: "root_id",
		}
		workloadName = "work"

		puller = &mockPuller{
			image: mockImage{ref: filter.Image, digest: "sha256:e454cab754cf9234e8b41d7c5e30f53a4c125d7d9443cb3ef2b2eb1c4bd1ec14"},
		}
		cancel     = func() {}
		deployment *appsv1.Deployment
	)
	BeforeEach(func() {
		cfg := aptest.MustConfig("")
		kube = kubernetes.NewForConfigOrDie(cfg)

		ns = "istio-provider-test-" + randutils.RandString(4)
		err := kubeutils.CreateNamespacesInParallel(kube, ns)
		Expect(err).NotTo(HaveOccurred())
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		mgr := aptest.ManagerWithOpts(ctx, cfg, manager.Options{
			Namespace: ns,
		})

		client = ezkube.NewEnsurer(ezkube.NewRestClient(mgr))

		cache = Cache{
			Namespace: ns,
			Name:      "cache-name",
		}

		deployment, err = kube.AppsV1().Deployments(ns).Create(makeDeployment(workloadName, ns))
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		cancel()
		if kube != nil {
			kubeutils.DeleteNamespacesInParallelBlocking(kube, ns)
		}
	})
	It("creates the EnvoyFilter", func() {
		workload := Workload{
			Labels:    deployment.Labels,
			Namespace: ns,
			Kind:      WorkloadTypeDeployment,
		}

		callbackCalled := false

		p, err := NewProvider(
			context.TODO(),
			kube,
			client,
			puller,
			workload,
			cache,
			nil,
			func(workloadMeta metav1.ObjectMeta, err error) {
				// test callback is called
				callbackCalled = true
			},
			"",
			0,
		)
		Expect(err).NotTo(HaveOccurred())

		err = p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep, err := kube.AppsV1().Deployments(workload.Namespace).Get(deployment.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		ef := &istiov1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      istioEnvoyFilterName(deployment.Name, filter.Id),
			},
		}
		err = client.Get(context.TODO(), ef)
		Expect(err).NotTo(HaveOccurred())

		Expect(ef.Spec.WorkloadSelector).To(Equal(&v1alpha3.WorkloadSelector{
			Labels: dep.Spec.Template.Labels,
		}))
		Expect(ef.Spec.ConfigPatches).To(HaveLen(1))

		Expect(callbackCalled).To(BeTrue())
	})
	It("given empty workload labels, annotates all workloads in the namespace and creates a generic EnvoyFilter", func() {
		workload := Workload{
			//all workloads
			Namespace: ns,
			Kind:      WorkloadTypeDeployment,
		}

		p := &Provider{
			Ctx:        context.TODO(),
			KubeClient: kube,
			Client:     client,
			Puller:     puller,
			Workload:   workload,
			Cache:      cache,
		}
		dep1, err := kube.AppsV1().Deployments(workload.Namespace).Create(makeDeployment("one", ns))
		Expect(err).NotTo(HaveOccurred())

		dep2, err := kube.AppsV1().Deployments(workload.Namespace).Create(makeDeployment("two", ns))
		Expect(err).NotTo(HaveOccurred())

		err = p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep1, err = kube.AppsV1().Deployments(workload.Namespace).Get(dep1.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		dep2, err = kube.AppsV1().Deployments(workload.Namespace).Get(dep2.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		ef1 := &istiov1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      istioEnvoyFilterName(dep1.Name, filter.Id),
			},
		}
		err = client.Get(context.TODO(), ef1)
		Expect(err).NotTo(HaveOccurred())

		Expect(ef1.Spec.WorkloadSelector.Labels).To(Equal(dep1.Spec.Template.Labels))
		Expect(ef1.Spec.ConfigPatches).To(HaveLen(1))

		ef2 := &istiov1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      istioEnvoyFilterName(dep2.Name, filter.Id),
			},
		}
		err = client.Get(context.TODO(), ef2)
		Expect(err).NotTo(HaveOccurred())

		Expect(ef2.Spec.WorkloadSelector.Labels).To(Equal(dep2.Spec.Template.Labels))
		Expect(ef2.Spec.ConfigPatches).To(HaveLen(1))
	})

	// note: this test assumes istio 1.5 installed to cluster
	It("returns an error when the image abi version does not support the istio version", func() {
		workload := Workload{
			//all workloads
			Namespace: ns,
			Kind:      WorkloadTypeDeployment,
		}
		resolver, _ := resolver.NewResolver("", "", false, false)
		puller := pull.NewPuller(resolver)
		client := mock_ezkube.NewMockEnsurer(gomock.NewController(GinkgoT()))

		p := &Provider{
			Ctx:        context.TODO(),
			KubeClient: kube,
			Client:     client,
			Puller:     puller,
			Workload:   workload,
			Cache:      cache,
		}
		glooImage := consts.HubDomain + "/ilackarms/gloo-test:1.3.3-0"
		err := p.ApplyFilter(&wasmev1.FilterSpec{
			Id:     "incompatible-filter",
			Image:  glooImage,
			Config: "{}",
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("image " + glooImage + " not supported by istio version"))

		client.EXPECT().Ensure(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err = p.ApplyFilter(&wasmev1.FilterSpec{
			Id:     "compatible-filter",
			Image:  test.IstioAssemblyScriptImage,
			Config: "{}",
		})
		Expect(err).NotTo(HaveOccurred())
	})
})

func makeDeployment(workloadName, ns string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": workloadName},
			},
			Template: kubev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": workloadName},
				},
				Spec: kubev1.PodSpec{
					Containers: []kubev1.Container{{
						Name:  "http-echo",
						Image: "hashicorp/http-echo",
						Args:  []string{fmt.Sprintf("-text=hi")},
						Ports: []kubev1.ContainerPort{{
							Name:          "http",
							ContainerPort: 5678,
						}},
					}},
					// important, otherwise termination lasts 30 seconds!
					TerminationGracePeriodSeconds: pointerToInt64(0),
				},
			},
		},
	}
}

type mockPuller struct {
	image mockImage
}

func (p *mockPuller) Pull(ctx context.Context, ref string) (pull.Image, error) {
	return &p.image, nil
}

type mockImage struct {
	ref    string
	digest string
}

func (m *mockImage) Ref() string {
	return m.ref
}

func (m *mockImage) Descriptor() (v1.Descriptor, error) {
	return v1.Descriptor{
		Digest: digest.Digest(m.digest),
	}, nil
}

func (m *mockImage) FetchFilter(ctx context.Context) (model.Filter, error) {
	panic("implement me")
}

func (m *mockImage) FetchConfig(ctx context.Context) (*config.Runtime, error) {
	return &config.Runtime{}, nil
}

func pointerToInt64(value int64) *int64 {
	return &value
}
