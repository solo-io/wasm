package istio

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/deploy"
	"istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	istiofake "istio.io/client-go/pkg/clientset/versioned/fake"
	appsv1 "k8s.io/api/apps/v1"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("IstioProvider", func() {
	var (
		kube  kubernetes.Interface
		istio versioned.Interface

		cache = Cache{
			Namespace: "cache-ns",
			Name:      "cache-name",
		}
		filter = &deploy.Filter{
			ID:     "filter-id",
			Config: `{"filter":"config"}`,
			Image:  "filter/image:v1",
			RootID: "root_id",
		}

		puller = &mockPuller{digest: "sha256:e454cab754cf9234e8b41d7c5e30f53a4c125d7d9443cb3ef2b2eb1c4bd1ec14"}
	)
	BeforeEach(func() {
		kube = fake.NewSimpleClientset()
		istio = istiofake.NewSimpleClientset()

		_, _ = kube.CoreV1().ConfigMaps(cache.Namespace).Create(&kubev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cache.Namespace,
				Name:      cache.Name,
			},
		})
	})
	It("annotates the workload and creates the EnvoyFilter", func() {
		workload := Workload{
			Name:      "work",
			Namespace: "load",
			Type:      WorkloadTypeDeployment,
		}

		p := &Provider{
			Ctx:         context.TODO(),
			KubeClient:  kube,
			IstioClient: istio,
			Puller:      puller,
			Workload:    workload,
			Cache:       cache,
		}

		podLabels := map[string]string{"these_labels": "expected_on_envoyfilter"}

		_, _ = kube.AppsV1().Deployments(workload.Namespace).Create(&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      workload.Name,
			},
			Spec: appsv1.DeploymentSpec{
				Template: kubev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: podLabels,
					},
					Spec: kubev1.PodSpec{
						Containers: []kubev1.Container{
							{
								Ports: []kubev1.ContainerPort{
									{
										ContainerPort: 123,
									},
									{
										ContainerPort: 456,
									},
								},
							},
							{
								Ports: []kubev1.ContainerPort{
									{
										ContainerPort: 789,
									},
								},
							},
						},
					},
				},
			},
		})

		err := p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep, err := kube.AppsV1().Deployments(workload.Namespace).Get(workload.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations([]uint32{123, 456, 789})))

		cacheConfig, err := kube.CoreV1().ConfigMaps(cache.Namespace).Get(cache.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(cacheConfig.Data).To(Equal(map[string]string{"images": filter.Image}))

		ef, err := istio.NetworkingV1alpha3().EnvoyFilters(workload.Namespace).Get(istioEnvoyFilterName(workload.Name, filter.ID), metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(ef.Spec.WorkloadSelector).To(Equal(&v1alpha3.WorkloadSelector{
			Labels: podLabels,
		}))
		Expect(ef.Spec.ConfigPatches).To(HaveLen(3))
		Expect(ef.Spec.ConfigPatches[0].Match.GetListener().PortNumber).To(Equal(uint32(123)))
		Expect(ef.Spec.ConfigPatches[1].Match.GetListener().PortNumber).To(Equal(uint32(456)))
		Expect(ef.Spec.ConfigPatches[2].Match.GetListener().PortNumber).To(Equal(uint32(789)))
	})
	It("given an empty workload name, annotates all workloads in the namespace and creates a generic EnvoyFilter", func() {
		workload := Workload{
			Name:      "", //all namespaces
			Namespace: "load",
			Type:      WorkloadTypeDeployment,
		}

		p := &Provider{
			Ctx:         context.TODO(),
			KubeClient:  kube,
			IstioClient: istio,
			Puller:      puller,
			Workload:    workload,
			Cache:       cache,
		}

		makeSpec := func(labels map[string]string) appsv1.DeploymentSpec {
			return appsv1.DeploymentSpec{
				Template: kubev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
					Spec: kubev1.PodSpec{
						Containers: []kubev1.Container{
							{
								Ports: []kubev1.ContainerPort{
									{
										ContainerPort: 123,
									},
									{
										ContainerPort: 456,
									},
								},
							},
							{
								Ports: []kubev1.ContainerPort{
									{
										ContainerPort: 789,
									},
								},
							},
						},
					},
				},
			}
		}

		labels1 := map[string]string{"hi": "there"}
		dep1, _ := kube.AppsV1().Deployments(workload.Namespace).Create(&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      "one",
			},
			Spec: makeSpec(labels1),
		})

		labels2 := map[string]string{"good": "bye"}
		dep2, _ := kube.AppsV1().Deployments(workload.Namespace).Create(&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      "two",
			},
			Spec: makeSpec(labels2),
		})

		err := p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep1, err = kube.AppsV1().Deployments(workload.Namespace).Get(dep1.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep1.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations([]uint32{123, 456, 789})))

		dep2, err = kube.AppsV1().Deployments(workload.Namespace).Get(dep2.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep1.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations([]uint32{123, 456, 789})))

		ef1, err := istio.NetworkingV1alpha3().EnvoyFilters(workload.Namespace).Get(istioEnvoyFilterName(dep1.Name, filter.ID), metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(ef1.Spec.WorkloadSelector.Labels).To(Equal(labels1))
		Expect(ef1.Spec.ConfigPatches).To(HaveLen(3))

		ef2, err := istio.NetworkingV1alpha3().EnvoyFilters(workload.Namespace).Get(istioEnvoyFilterName(dep2.Name, filter.ID), metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(ef2.Spec.WorkloadSelector.Labels).To(Equal(labels2))
		Expect(ef2.Spec.ConfigPatches).To(HaveLen(3))
	})
})

type mockPuller struct {
	digest string
}

func (m *mockPuller) PullCodeDescriptor(ctx context.Context, ref string) (v1.Descriptor, error) {
	return v1.Descriptor{
		Digest: digest.Digest(m.digest),
	}, nil
}
