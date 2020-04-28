package gloo_test

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"time"

	skutil "github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/wasme/pkg/util"
	"github.com/solo-io/wasme/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = AfterSuite(func() {
	// delete gloo-system-test to make room for other things in the cluster
	if err := skutil.Kubectl(nil, "delete", "ns", "gloo-system-test"); err != nil {
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
				"http://localhost:8001/api/v1/namespaces/gloo-system-test/services/gateway-proxy:http/proxy/api/pets")

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
