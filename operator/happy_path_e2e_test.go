package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/codegen/model"
	"github.com/solo-io/autopilot/codegen/render"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/autopilot/cli/pkg/utils"
	"github.com/solo-io/autopilot/codegen/util"
)

func runMake(target string) error {
	cmd := exec.Command("make", "-C", filepath.Dir(util.GoModPath()), target)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func applyFile(file, ns string) error {
	return withManifest(file, ns, utils.KubectlApply)
}

func deleteFile(file, ns string) error {
	return withManifest(file, ns, utils.KubectlDelete)
}

func withManifest(file, ns string, fn func(manifest []byte, extraArgs ...string) error) error {
	path := filepath.Join(util.MustGetThisDir(), file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	extraArgs := []string{}
	if ns != "" {
		extraArgs = []string{"-n", ns}
	}
	return fn(b, extraArgs...)
}

func generateCrdExample() error {
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

	filename := filepath.Join(util.MustGetThisDir(), "test_filter.yaml")
	if err := ioutil.WriteFile(filename, []byte(filterDeploymentFile[0].Content), 0644); err != nil {
		return err
	}

	return nil
}

var ns = "bookinfo"

var _ = BeforeSuite(func() {
	err := runMake("manifest-gen")
	Expect(err).NotTo(HaveOccurred())

	err = generateCrdExample()
	Expect(err).NotTo(HaveOccurred())

	// ensure no collision between tests
	err = waitNamespaceTerminated(ns, time.Minute)
	Expect(err).NotTo(HaveOccurred())

	utils.Kubectl(nil, "create", "ns", ns)

	err = utils.Kubectl(nil, "label", "namespace", ns, "istio-injection=enabled", "--overwrite")
	Expect(err).NotTo(HaveOccurred())

	err = applyFile("install/wasme/crds/wasme.io_v1_crds.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	err = applyFile("install/wasme-default.yaml", "")
	Expect(err).NotTo(HaveOccurred())

	err = applyFile("bookinfo.yaml", ns)
	Expect(err).NotTo(HaveOccurred())

	err = waitDeploymentReady("productpage", "bookinfo", time.Minute*2)
	Expect(err).NotTo(HaveOccurred())
})
var _ = AfterSuite(func() {
	deleteFile("bookinfo.yaml", ns)
	deleteFile("install/wasme-default.yaml", "")
	utils.Kubectl(nil, "delete", "ns", ns)
})

var _ = Describe("AutopilotGenerate", func() {
	It("runs the wasme operator", func() {

		err := applyFile("test_filter.yaml", ns)
		Expect(err).NotTo(HaveOccurred())

		testRequest := func() (string, error) {
			return utils.KubectlOut(nil,
				"exec",
				"-n", ns,
				"deploy/productpage-v1",
				"-c", "istio-proxy", "--",
				"curl", "-v", "http://details."+ns+":9080/details/123")
		}

		// expect header in response
		Eventually(testRequest, time.Minute*5).Should(ContainSubstring("hello: world"))

		err = deleteFile("test_filter.yaml", ns)
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
