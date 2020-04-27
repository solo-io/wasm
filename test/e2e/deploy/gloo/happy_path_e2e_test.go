package gloo_test

import (
	"bytes"
	"log"
	"strings"
	"time"

	skutil "github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/wasme/pkg/util"

	"github.com/solo-io/wasme/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

// use a namespace as a "cluster lock"
var ns = "gloo-e2e-happy-path-test-lock"

var _ = BeforeSuite(func() {
	// ensure no collision between tests
	err := waitNamespaceTerminated(ns, time.Minute)
	Expect(err).NotTo(HaveOccurred())

	err = skutil.Kubectl(nil, "create", "ns", ns)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if err := skutil.Kubectl(nil, "delete", "ns", ns); err != nil {
		log.Printf("failed deleting ns: %v", err)
	}
})

// Test Order matters here.
// Do not randomize ginkgo specs when running, if the build & push test is enabled
var _ = Describe("wasme deploy gloo", func() {
	It("Deploys the filter via Gloo", func() {
		imageName := test.GetImageTag()

		err := test.WasmeCli("deploy", "gloo", imageName, "--id", "myfilter", "--config", "world")
		Expect(err).NotTo(HaveOccurred())

		gatewayAddr, err := util.ExecOutput(nil, "kubectl", "get", "svc", "-n", "gloo-system", "gateway-proxy", "-o", "jsonpath={.status.loadBalancer.ingress[*].ip}")
		Expect(err).NotTo(HaveOccurred())

		testRequest := func() (string, error) {
			b := &bytes.Buffer{}
			err := util.ExecCmd(
				b,
				b,
				nil,
				"curl",
				"-v",
				"http://"+gatewayAddr+"/api/pets")

			out := b.String()

			return out, errors.Wrapf(err, out)
		}

		// expect header in response
		Eventually(testRequest, time.Minute*5).Should(ContainSubstring("hello: world"))

		err = test.WasmeCli("undeploy", "gloo", "--id", "myfilter")
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
			out, err := skutil.KubectlOut(nil, "get", "pod", "-n", namespace, "-l", "app="+name)
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
			_, err := skutil.KubectlOut(nil, "get", "namespace", namespace)
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
