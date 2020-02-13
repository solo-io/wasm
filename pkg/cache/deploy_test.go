package cache_test

import (
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/randutils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/solo-io/wasme/pkg/cache"
)

var _ = Describe("Deploy", func() {
	var (
		kube kubernetes.Interface
		// switch to running the test vs a real kube cluster
		useRealKube = os.Getenv("USE_REAL_KUBE") != ""

		cacheNamespace = "wasme-cache-test-" + randutils.RandString(4)
	)

	BeforeEach(func() {
		if useRealKube {
			cfg, err := kubeutils.GetConfig("", "")
			Expect(err).NotTo(HaveOccurred())

			kube, err = kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		} else {
			kube = fake.NewSimpleClientset()
		}
	})
	AfterEach(func() {
		kube.AppsV1().DaemonSets(cacheNamespace).Delete(CacheName, nil)
		kube.CoreV1().ConfigMaps(cacheNamespace).Delete(CacheName, nil)
		kube.CoreV1().Namespaces().Delete(cacheNamespace, nil)
	})
	It("creates the cache namespace, configmap, and daemonset", func() {

		deployer := NewDeployer(kube, cacheNamespace, "", "", "", nil, corev1.PullAlways)

		err := deployer.EnsureCache()
		Expect(err).NotTo(HaveOccurred())

		_, err = kube.CoreV1().Namespaces().Get(cacheNamespace, v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		_, err = kube.CoreV1().ConfigMaps(cacheNamespace).Get(CacheName, v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		_, err = kube.AppsV1().DaemonSets(cacheNamespace).Get(CacheName, v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		if !useRealKube {
			return
		}

		// multiple runs should not error
		err = deployer.EnsureCache()
		Expect(err).NotTo(HaveOccurred())

		// eventually pods should be ready
		Eventually(func() (int32, error) {
			cacheDaemonSet, err := kube.AppsV1().DaemonSets(cacheNamespace).Get(CacheName, v1.GetOptions{})
			if err != nil {
				return 0, err
			}
			return cacheDaemonSet.Status.NumberReady, nil
		}, time.Second*30).Should(Equal(int32(1)))

	})
})
