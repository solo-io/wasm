package cache_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	testutils "github.com/solo-io/wasme/test"

	corev1 "k8s.io/api/core/v1"

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

		// used for pushing images when USE_REAL_KUBE is true
		operatorImage = func() string {
			if gcloudProject := os.Getenv("GCLOUD_PROJECT_ID"); gcloudProject != "" {
				return fmt.Sprintf("gcr.io/%v/wasme", gcloudProject)
			}
			return "quay.io/solo-io/wasme"
		}()
	)

	BeforeEach(func() {
		if useRealKube {
			err := testutils.RunMake("wasme-image", func(cmd *exec.Cmd) {
				cmd.Args = append(cmd.Args, "OPERATOR_IMAGE="+operatorImage)
				cmd.Args = append(cmd.Args, "VERSION="+cacheNamespace)
			})
			Expect(err).NotTo(HaveOccurred())

			err = testutils.RunMake("wasme-image-push", func(cmd *exec.Cmd) {
				cmd.Args = append(cmd.Args, "OPERATOR_IMAGE="+operatorImage)
				cmd.Args = append(cmd.Args, "VERSION="+cacheNamespace)
			})
			Expect(err).NotTo(HaveOccurred())

			cfg, err := kubeutils.GetConfig("", "")
			Expect(err).NotTo(HaveOccurred())

			kube, err = kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
		} else {
			kube = fake.NewSimpleClientset()
		}
	})
	AfterEach(func() {
		kube.AppsV1().Deployments(cacheNamespace).Delete(CacheName, nil)
		kube.CoreV1().ConfigMaps(cacheNamespace).Delete(CacheName, nil)
		kube.CoreV1().Namespaces().Delete(cacheNamespace, nil)
	})
	It("creates the cache namespace, and deployment", func() {

		deployer := NewDeployer(kube, cacheNamespace, "", operatorImage, cacheNamespace, nil, corev1.PullAlways)

		err := deployer.EnsureCache()
		Expect(err).NotTo(HaveOccurred())

		_, err = kube.CoreV1().Namespaces().Get(cacheNamespace, v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		_, err = kube.AppsV1().Deployments(cacheNamespace).Get(CacheName, v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		if !useRealKube {
			return
		}

		// multiple runs should not error
		err = deployer.EnsureCache()
		Expect(err).NotTo(HaveOccurred())

		// eventually pods should be ready
		Eventually(func() (int32, error) {
			cacheDeployments, err := kube.AppsV1().Deployments(cacheNamespace).Get(CacheName, v1.GetOptions{})
			if err != nil {
				return 0, err
			}
			return cacheDeployments.Status.ReadyReplicas, nil
		}, time.Second*30).Should(Equal(int32(1)))

	})
})
