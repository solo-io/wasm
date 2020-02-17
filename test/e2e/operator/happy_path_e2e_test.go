package operator_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/solo-io/autopilot/codegen/util"

	testdefaults "github.com/solo-io/wasme/pkg/consts/test"
	"github.com/solo-io/wasme/test"

	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/codegen/model"
	"github.com/solo-io/autopilot/codegen/render"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/autopilot/cli/pkg/utils"
)

func generateCrdExample(filename, image string) error {
	if image == "" {
		image = testdefaults.IstioAssemblyScriptImage
	}

	filterDeployment := &v1.FilterDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FilterDeployment",
			APIVersion: "wasme.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myfilter",
			Namespace: "bookinfo",
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

	utils.Kubectl(nil, "create", "ns", ns)

	err = utils.Kubectl(nil, "label", "namespace", ns, "istio-injection=enabled", "--overwrite")
	Expect(err).NotTo(HaveOccurred())

	err = test.ApplyFile("operator/install/wasme/crds/wasme.io_v1_crds.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	err = test.ApplyFile("operator/install/wasme-default.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	err = test.ApplyFile("test/e2e/operator/bookinfo.yaml", ns)
	Expect(err).NotTo(HaveOccurred())

	err = waitDeploymentReady("productpage", "bookinfo", time.Minute*2)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if err := test.DeleteFile("test/e2e/operator/bookinfo.yaml", ns); err != nil {
		log.Printf("failed deleting file: %v", err)
	}
	if err := test.DeleteFile("operator/install/wasme-default.yaml", ""); err != nil {
		log.Printf("failed deleting file: %v", err)
	}
	if err := utils.Kubectl(nil, "delete", "ns", ns); err != nil {
		log.Printf("failed deleting ns: %v", err)
	}
})

// Test Order matters here.
// Do not randomize ginkgo specs when running, if the build & push test is enabled
var _ = Describe("AutopilotGenerate", func() {
	It("runs the wasme operator", func() {
		filterFile := "test/e2e/operator/test_filter.yaml"

		err := generateCrdExample(filterFile, os.Getenv("FILTER_IMAGE_TAG"))
		Expect(err).NotTo(HaveOccurred())

		err = test.ApplyFile(filterFile, ns)
		Expect(err).NotTo(HaveOccurred())

		testRequest := func() (string, error) {
			out, err := utils.KubectlOut(nil,
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
			out, err := utils.KubectlOut(nil, "get", "pod", "-n", namespace, "-l", "app="+name)
			if err != nil {
				return err
			}
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
			_, err := utils.KubectlOut(nil, "get", "namespace", namespace)
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
