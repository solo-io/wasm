package build_test

import (
	"os"
	"path/filepath"

	"github.com/solo-io/wasme/pkg/consts"

	"github.com/solo-io/autopilot/codegen/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/wasme/test"
)

var _ = Describe("Build", func() {
	It("builds and pushes the image", func() {
		imageName := test.GetImageTag()
		username := test.GetEnv("WASME_LOGIN_USERNAME")
		password := test.GetEnv("WASME_LOGIN_PASSWORD")

		err := test.RunMake("generated-code")
		Expect(err).NotTo(HaveOccurred())

		err = test.WasmeCliSplit("login -u " + username + " -p " + password + " -s " + consts.HubDomain)
		Expect(err).NotTo(HaveOccurred())

		err = test.WasmeCliSplit("init test-filter --platform istio --platform-version 1.5.x --language assemblyscript")
		Expect(err).NotTo(HaveOccurred())

		if precompiledFilter := os.Getenv("PRECOMPILED_FILTER_PATH"); precompiledFilter != "" {
			err = test.WasmeCliSplit("build precompiled -t " + imageName + " test-filter " + filepath.Dir(util.GoModPath()) + "/" + precompiledFilter)
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
				// TODO: remove login info when NPM repo is published
			)
			Expect(err).NotTo(HaveOccurred())
		}

		err = test.WasmeCliSplit("push " + imageName)
		Expect(err).NotTo(HaveOccurred())
	})
})
