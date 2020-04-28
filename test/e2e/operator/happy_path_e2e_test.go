package operator_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/solo-io/wasme/pkg/cache"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/solo-io/skv2/codegen/util"

	"github.com/solo-io/wasme/test"

	"github.com/pkg/errors"
	"github.com/solo-io/skv2/codegen/model"
	"github.com/solo-io/skv2/codegen/render"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var filterDeploymentName = "myfilter"

func generateCrdExample(filename, image, ns string) error {
	filterDeployment := &v1.FilterDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FilterDeployment",
			APIVersion: "wasme.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      filterDeploymentName,
			Namespace: ns,
		},
		Spec: v1.FilterDeploymentSpec{
			Filter: &v1.FilterSpec{
				Image:  image,
				Config: `world`,
			},
			Deployment: &v1.DeploymentSpec{
				DeploymentType: &v1.DeploymentSpec_Istio{Istio: &v1.IstioDeploymentSpec{
					Kind: "Deployment",
				}},
			},
		},
	}

	// hack to write the file as yaml
	filterDeploymentFile, err := render.ManifestsRenderer{
		AppName: "wasme-test-app",
		ResourceFuncs: map[render.OutFile]render.MakeResourceFunc{
			render.OutFile{}: func(group render.Group) []metav1.Object {
				return []metav1.Object{filterDeployment}
			},
		},
	}.RenderManifests(model.Group{RenderManifests: true})
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(filepath.Dir(util.GoModPath()), filename), []byte(filterDeploymentFile[0].Content), 0644); err != nil {
		return err
	}

	return nil
}

var ns = "bookinfo"

var _ = BeforeSuite(func() {
	err := test.RunMake("manifest-gen")
	Expect(err).NotTo(HaveOccurred())

	// ensure no collision between tests
	err = waitNamespaceTerminated(ns, time.Minute)
	Expect(err).NotTo(HaveOccurred())

	util.Kubectl(nil, "create", "ns", ns)

	err = util.Kubectl(nil, "label", "namespace", ns, "istio-injection=enabled", "--overwrite")
	Expect(err).NotTo(HaveOccurred())

	err = test.ApplyFile("operator/install/wasme/crds/wasme.io_v1_crds.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	err = test.ApplyFile("operator/install/wasme-default.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	patchCacheDaemonSet()

	err = test.ApplyFile("test/e2e/operator/bookinfo.yaml", ns)
	Expect(err).NotTo(HaveOccurred())

	err = waitDeploymentReady("productpage", ns, time.Minute*5)
	Expect(err).NotTo(HaveOccurred())
})

// need to patch the cache daemonset to use the --clear-cache flag, to ensure
// our cache starts fresh every test
func patchCacheDaemonSet() {
	cfg, err := config.GetConfig()
	Expect(err).NotTo(HaveOccurred())

	kube, err := client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())

	ds := &appsv1.DaemonSet{}
	err = kube.Get(context.TODO(), client.ObjectKey{Name: cache.CacheName, Namespace: cache.CacheNamespace}, ds)
	Expect(err).NotTo(HaveOccurred())

	args := ds.Spec.Template.Spec.Containers[0].Args
	args = append(args, "--clear-cache")
	ds.Spec.Template.Spec.Containers[0].Args = args

	err = kube.Update(context.TODO(), ds)
	Expect(err).NotTo(HaveOccurred())
}

var _ = AfterSuite(func() {
	if err := test.DeleteFile("test/e2e/operator/bookinfo.yaml", ns); err != nil {
		log.Printf("failed deleting file: %v", err)
	}
	if err := test.DeleteFile("operator/install/wasme-default.yaml", ""); err != nil {
		log.Printf("failed deleting file: %v", err)
	}
	if err := util.Kubectl(nil, "delete", "ns", ns); err != nil {
		log.Printf("failed deleting ns: %v", err)
	}
})

// Test Order matters here.
// Do not randomize ginkgo specs when running, if the build & push test is enabled
var _ = Describe("skv2Generate", func() {
	It("runs the wasme operator", func() {
		filterFile := "test/e2e/operator/test_filter.yaml"

		err := generateCrdExample(filterFile, test.GetImageTag(), ns)
		Expect(err).NotTo(HaveOccurred())

		err = test.ApplyFile(filterFile, ns)
		Expect(err).NotTo(HaveOccurred())

		testRequest := func() (string, error) {
			out, err := util.KubectlOut(nil,
				"exec",
				"-n", ns,
				"deploy/productpage-v1",
				"-c", "istio-proxy", "--",
				"curl", "-v", "http://details."+ns+":9080/details/123")

			log.Printf("output: %v", out)
			log.Printf("err: %v", err)
			return out, err
		}

		// expect header in response
		Eventually(testRequest, time.Minute*5).Should(ContainSubstring("hello: world"))

		// ensure filter deployment status is up to date
		cfg, err := config.GetConfig()
		Expect(err).NotTo(HaveOccurred())

		err = v1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		kube, err := client.New(cfg, client.Options{})
		Expect(err).NotTo(HaveOccurred())

		fd := &v1.FilterDeployment{}
		Eventually(func() (int64, error) {
			err := kube.Get(context.TODO(), client.ObjectKey{Name: filterDeploymentName, Namespace: ns}, fd)
			if err != nil {
				return 0, err
			}
			return fd.Status.ObservedGeneration, nil
		}).Should(Equal(int64(1)))

		err = test.DeleteFile(filterFile, ns)
		Expect(err).NotTo(HaveOccurred())

		// expect header not in response
		Eventually(testRequest, time.Minute*3).ShouldNot(ContainSubstring("hello: world"))

	})
})

func waitDeploymentReady(name, namespace string, timeout time.Duration) error {
	timedOut := time.After(timeout)
	for {
		select {
		case <-timedOut:
			return errors.Errorf("timed out after %s", timeout)
		default:
			out, err := util.KubectlOut(nil, "get", "pod", "-n", namespace, "-l", "app="+name)
			if err != nil {
				return err
			}
			fmt.Println(GinkgoWriter, "waiting for deployment: pod status", string(out))
			if strings.Contains(out, "Running") && strings.Contains(out, "2/2") {
				return nil
			}
			time.Sleep(time.Second * 2)
		}
	}
}

func waitNamespaceTerminated(namespace string, timeout time.Duration) error {
	timedOut := time.After(timeout)
	for {
		select {
		case <-timedOut:
			return errors.Errorf("timed out after %s", timeout)
		default:
			_, err := util.KubectlOut(nil, "get", "namespace", namespace)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					return nil
				}
				return err
			}
			time.Sleep(time.Second * 2)
		}
	}
}
