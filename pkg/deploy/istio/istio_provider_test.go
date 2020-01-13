package istio

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/solo-io/autopilot/cli/pkg/utils"
	"github.com/solo-io/autopilot/codegen/util"
	"github.com/solo-io/autopilot/pkg/ezkube"
	"github.com/solo-io/autopilot/test"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/randutils"
	wasmev1 "github.com/solo-io/wasme/operator/pkg/api/wasme.io/v1"
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

func applyCrds() error {
	// apply operator crd
	path := filepath.Join(util.GetModuleRoot(), "operator/install/kube/wasme.io_v1_crds.yaml")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return utils.KubectlApply(b)
}

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

		puller = &mockPuller{digest: "sha256:e454cab754cf9234e8b41d7c5e30f53a4c125d7d9443cb3ef2b2eb1c4bd1ec14"}

		cancel = func() {}
	)
	BeforeEach(func() {
		err := applyCrds()
		Expect(err).NotTo(HaveOccurred())

		cfg := test.MustConfig()
		kube = kubernetes.NewForConfigOrDie(cfg)

		ns = "istio-provider-test-" + randutils.RandString(4)
		err = kubeutils.CreateNamespacesInParallel(kube, ns)
		Expect(err).NotTo(HaveOccurred())

		mgr, c := test.ManagerWithOpts(cfg, manager.Options{
			Namespace:               ns,
			LeaderElection:          true,
			LeaderElectionNamespace: ns,
		})
		cancel = c

		client = ezkube.NewEnsurer(ezkube.NewRestClient(mgr))

		cache = Cache{
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

		_, err = kube.AppsV1().Deployments(ns).Create(makeDeployment(workloadName, ns))
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		cancel()
		if kube != nil {
			kubeutils.DeleteNamespacesInParallel(kube, ns)
		}
	})
	It("annotates the workload and creates the EnvoyFilter", func() {
		workload := Workload{
			Name:      workloadName,
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
		)
		Expect(err).NotTo(HaveOccurred())

		err = p.ApplyFilter(filter)
		Expect(err).NotTo(HaveOccurred())

		dep, err := kube.AppsV1().Deployments(workload.Namespace).Get(workload.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(dep.Spec.Template.Annotations).To(Equal(requiredSidecarAnnotations()))

		cacheConfig, err := kube.CoreV1().ConfigMaps(cache.Namespace).Get(cache.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		Expect(cacheConfig.Data).To(Equal(map[string]string{"images": filter.Image}))

		ef := &istiov1alpha3.EnvoyFilter{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: workload.Namespace,
				Name:      istioEnvoyFilterName(workload.Name, filter.Id),
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
	It("given an empty workload name, annotates all workloads in the namespace and creates a generic EnvoyFilter", func() {
		workload := Workload{
			Name:      "", //all workloads
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
})

type mockPuller struct {
	digest string
}

func (m *mockPuller) PullCodeDescriptor(ctx context.Context, ref string) (v1.Descriptor, error) {
	return v1.Descriptor{
		Digest: digest.Digest(m.digest),
	}, nil
}

func pointerToInt64(value int64) *int64 {
	return &value
}
