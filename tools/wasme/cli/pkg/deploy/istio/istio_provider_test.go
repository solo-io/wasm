package istio_test

import (
	"context"
	"fmt"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/deploy/istio"
	"github.com/solo-io/wasm/tools/wasme/pkg/config"
	"github.com/solo-io/wasm/tools/wasme/pkg/consts"
	"github.com/solo-io/wasm/tools/wasme/pkg/model"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"

	"github.com/golang/mock/gomock"
	mock_ezkube "github.com/solo-io/skv2/pkg/ezkube/mocks"

	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/randutils"
	"github.com/solo-io/skv2/pkg/ezkube"

	aptest "github.com/solo-io/skv2/test"
	wasmev1 "github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1"
	testutils "github.com/solo-io/wasm/tools/wasme/cli/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/api/networking/v1alpha3"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("IstioProvider", func() {
	// TODO (shane) remove arbitrary change to trigger CI

	var (
		kube   kubernetes.Interface
		client ezkube.Ensurer
		ns     string
		cache  istio.Cache
		filter = &wasmev1.FilterSpec{
			Id:     "filter-id",
			Config: nil,
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

		ctx, c := context.WithCancel(context.Background())
		mgr := aptest.ManagerWithOpts(ctx, cfg, manager.Options{
			Namespace: ns,
		})
		cancel = c

		client = ezkube.NewEnsurer(ezkube.NewRestClient(mgr))

		cache = istio.Cache{
			Namespace: ns,
			Name:      "cache-name",
		}
		_, err = kube.CoreV1().ConfigMaps(ns).Create(&kubev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cache.Namespace,
				Name:      cache.Name,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		deployment, err = kube.AppsV1().Deployments(ns).Create(makeDeployment(workloadName, ns))
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		cancel()
		if kube != nil {
			kubeutils.DeleteNamespacesInParallelBlocking(kube, ns)
		}
	})
	It("annotates the workload and creates the EnvoyFilter", func() {
		workload := istio.Workload{
			Labels:    deployment.Labels,
			Namespace: ns,
			Kind:      istio.WorkloadTypeDeployment,
		}

		callbackCalled := false

		p, err := istio.NewProvider(
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
			false,
		)
		Expect(err).NotTo(HaveOccurred())

		err = p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep, err := kube.AppsV1().Deployments(workload.Namespace).Get(deployment.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations()))

		cacheConfig, err := kube.CoreV1().ConfigMaps(cache.Namespace).Get(cache.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(cacheConfig.Data).To(Equal(map[string]string{"images": filter.Image}))

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
		workload := istio.Workload{
			//all workloads
			Namespace: ns,
			Kind:      istio.WorkloadTypeDeployment,
		}

		p := &istio.Provider{
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

		Expect(dep1.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations()))

		dep2, err = kube.AppsV1().Deployments(workload.Namespace).Get(dep2.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep1.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations()))

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

	It("create an Envoy filter for outbound traffic", func() {
		workload := istio.Workload{
			//all workloads
			Namespace: ns,
			Kind:      istio.WorkloadTypeDeployment,
		}

		p := &istio.Provider{
			Ctx:        context.TODO(),
			KubeClient: kube,
			Client:     client,
			Puller:     puller,
			Workload:   workload,
			Cache:      cache,
		}

		obfilter := &wasmev1.FilterSpec{
			Id:           "filter-id",
			Config:       nil,
			Image:        "filter/image:v1",
			RootID:       "root_id",
			PatchContext: "outbound",
		}
		err := p.ApplyFilter(obfilter)
		Expect(err).NotTo(HaveOccurred())

		ef := &istiov1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      istioEnvoyFilterName(deployment.Name, obfilter.Id),
			},
		}
		err = client.Get(context.TODO(), ef)
		Expect(err).NotTo(HaveOccurred())

		Expect(ef.Spec.ConfigPatches).To(HaveLen(1))
		Expect(ef.Spec.ConfigPatches[0].Match.Context).To(Equal(networkingv1alpha3.EnvoyFilter_SIDECAR_OUTBOUND))
	})

	// note: this test assumes istio 1.5 installed to cluster
	It("returns an error when the image abi version does not support the istio version", func() {
		workload := istio.Workload{
			//all workloads
			Namespace: ns,
			Kind:      istio.WorkloadTypeDeployment,
		}
		resolver, _ := resolver.NewResolver("", "", false, false)
		puller := pull.NewPuller(resolver)
		client := mock_ezkube.NewMockEnsurer(gomock.NewController(GinkgoT()))

		p := &istio.Provider{
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
			Config: nil,
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("image " + glooImage + " not supported by istio version"))

		client.EXPECT().Ensure(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		err = p.ApplyFilter(&wasmev1.FilterSpec{
			Id:     "compatible-filter",
			Image:  testutils.GetImageTagIstio(),
			Config: nil,
		})
		Expect(err).NotTo(HaveOccurred())
	})
	// note: this test assumes istio 1.5 installed to cluster
	It("returns no errors when the image abi version does not explicitly support the istio version, but --ignore-version-check is set", func() {
		workload := istio.Workload{
			//all workloads
			Namespace: ns,
			Kind:      istio.WorkloadTypeDeployment,
		}
		resolver, _ := resolver.NewResolver("", "", false, false)
		puller := pull.NewPuller(resolver)
		client := mock_ezkube.NewMockEnsurer(gomock.NewController(GinkgoT()))

		p := &istio.Provider{
			Ctx:                context.TODO(),
			KubeClient:         kube,
			Client:             client,
			Puller:             puller,
			Workload:           workload,
			Cache:              cache,
			IngoreVersionCheck: true,
		}

		client.EXPECT().Ensure(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		glooImage := consts.HubDomain + "/ilackarms/gloo-test:1.3.3-0"
		incompatibleFilter := &wasmev1.FilterSpec{
			Id:     "incompatible-filter",
			Image:  glooImage,
			Config: nil,
		}
		err := p.ApplyFilter(incompatibleFilter)
		Expect(err).NotTo(HaveOccurred())

		// Since this filter won't actually work (it's not compatible),
		// we need to remove it again so we're not messing up the cluster
		client.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1)
		p.RemoveFilter(incompatibleFilter)

		err = p.ApplyFilter(&wasmev1.FilterSpec{
			Id:     "compatible-filter",
			Image:  testutils.GetImageTagIstio(),
			Config: nil,
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

func istioEnvoyFilterName(workloadName, filterId string) string {
	return workloadName + "-" + filterId
}

// the sidecar annotations required on the pod
func requiredSidecarAnnotations() map[string]string {
	return map[string]string{
		"sidecar.istio.io/userVolume":      `[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache"}}]`,
		"sidecar.istio.io/userVolumeMount": `[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]`,
	}
}
