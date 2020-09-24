package build_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	skutil "github.com/solo-io/skv2/codegen/util"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/wasm/tools/wasme/cli/test"
	"github.com/solo-io/wasm/tools/wasme/pkg/consts"
)

var _ = Describe("Build", func() {

	AfterEach(func() {
		os.RemoveAll("test-filter")
	})

	It("builds and pushes the image", func() {
		imageName := test.GetBuildImageTag()
		username := os.Getenv("WASME_LOGIN_USERNAME")
		password := os.Getenv("WASME_LOGIN_PASSWORD")
		var err error
		if password != "" {
			err = test.WasmeCliSplit("login -u " + username + " -p " + password + " -s " + consts.HubDomain)
			Expect(err).NotTo(HaveOccurred())
		}

		err = test.WasmeCliSplit("init test-filter --platform istio --platform-version 1.5.x --language assemblyscript")
		Expect(err).NotTo(HaveOccurred())

		if precompiledFilter := os.Getenv("PRECOMPILED_FILTER_PATH"); precompiledFilter != "" {
			err = test.WasmeCliSplit("build precompiled -t " + imageName + " test-filter " + filepath.Dir(skutil.GoModPath()) + "/" + precompiledFilter)
			Expect(err).NotTo(HaveOccurred())
		} else {
			err = test.RunMake("builder-image")
			Expect(err).NotTo(HaveOccurred())

			err = test.WasmeCli(
				"build",
				"assemblyscript",
				// need to run with --tmp-dir=. in CI due to docker mount concerns
				"--tmp-dir=.",
				"-t="+imageName,
				"test-filter",
			)
			Expect(err).NotTo(HaveOccurred())
		}

		err = test.WasmeCliSplit("push " + imageName)
		Expect(err).NotTo(HaveOccurred())
	})

	ExpectRustExampleToWorkInIstioVersion := func(version string) {
		err := test.WasmeCliSplit("init test-filter --platform istio --platform-version " + version + ".x --language rust")
		Expect(err).NotTo(HaveOccurred())
		imagename := "testimage"
		envoyimage := "docker.io/istio/proxyv2:" + version + ".0"
		err = test.RunMake("builder-image")
		Expect(err).NotTo(HaveOccurred())
		err = test.WasmeCli(
			"build",
			"rust",
			// need to run with --tmp-dir=. in CI due to docker mount concerns
			"--tmp-dir=.",
			"-t="+imagename,
			"test-filter",
		)
		Expect(err).NotTo(HaveOccurred())

		go func() {
			defer GinkgoRecover()
			err := test.WasmeCli(
				"deploy",
				"envoy",
				imagename,
				"--envoy-image="+envoyimage,
				"--envoy-run-args=-l trace",
				"--id=myfilter",
			)
			Expect(err).NotTo(HaveOccurred())
		}()

		testRequest := func() (string, error) {
			b := &bytes.Buffer{}
			w := io.MultiWriter(b, GinkgoWriter)
			err := util.ExecCmd(
				w,
				w,
				nil,
				"curl",
				"-v",
				"http://localhost:8080/")
			out := b.String()
			return out, errors.Wrapf(err, out)
		}

		// expect header in response
		// note that header key is capital case as this goes through Kube api
		const addedHeader = "hello: world"
		Eventually(testRequest, 10*time.Second, time.Second).Should(ContainSubstring(addedHeader))

		util.Docker(GinkgoWriter, GinkgoWriter, nil, "stop", "myfilter")
	}

	Context("istio-rust", func() {
		It("builds a valid image", func() {
			By("istio 1.5")
			ExpectRustExampleToWorkInIstioVersion("1.5")
			os.RemoveAll("test-filter")
			By("istio 1.6")
			ExpectRustExampleToWorkInIstioVersion("1.6")
			os.RemoveAll("test-filter")
			By("istio 1.7")
			ExpectRustExampleToWorkInIstioVersion("1.7")
		})
	})
})
