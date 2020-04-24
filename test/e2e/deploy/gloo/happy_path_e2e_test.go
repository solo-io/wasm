package gloo_test

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/solo-io/wasme/pkg/util"

	"github.com/solo-io/wasme/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/cli/pkg/utils"
)

// use a namespace as a "cluster lock"
var ns = "gloo-e2e-happy-path-test-lock"

var _ = BeforeSuite(func() {
	// ensure no collision between tests
	err := waitNamespaceTerminated(ns, time.Minute)
	Expect(err).NotTo(HaveOccurred())

	err = utils.Kubectl(nil, "create", "ns", ns)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if err := utils.Kubectl(nil, "delete", "ns", ns); err != nil {
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		exec.CommandContext(ctx, "kubectl", "proxy").Start()

		testRequest := func() (string, error) {
			b := &bytes.Buffer{}
			err := util.ExecCmd(
				b,
				b,
				nil,
				"curl",
				"-v",
				"http://localhost:8001/api/v1/namespaces/gloo-system/services/gateway-proxy:http/proxy/api/pets")

			out := b.String()

			return out, errors.Wrapf(err, out)
		}

		// expect header in response
		// note that header key is capital case as this goes through Kube api
		const addedHeader = "Hello: world"
		Eventually(testRequest, time.Minute*5).Should(ContainSubstring(addedHeader))

		err = test.WasmeCli("undeploy", "gloo", "--id", "myfilter")
		Expect(err).NotTo(HaveOccurred())

		// expect header not in response
		Eventually(testRequest, time.Minute*3).ShouldNot(ContainSubstring(addedHeader))
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
