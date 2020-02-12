package build_test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/wasme/pkg/consts"

	"github.com/solo-io/autopilot/codegen/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/wasme/test"
)

func getEnv(env string) string {
	val := strings.TrimSpace(os.Getenv(env))
	if val == "" {
		Skip("Skipping build/push test. To enable, set " + env + " to the tag to use for the built/pushed image")
	}
	return val
}

var _ = Describe("Build", func() {
	It("builds and pushes the image", func() {
		imageName := getEnv("FILTER_IMAGE_TAG")
		username := getEnv("WASME_LOGIN_USERNAME")
		password := getEnv("WASME_LOGIN_PASSWORD")
		npmUsername := getEnv("NPM_LOGIN_USERNAME")
		npmPassword := getEnv("NPM_LOGIN_PASSWORD")
		npmEmail := getEnv("NPM_LOGIN_EMAIL")

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
				"-u="+npmUsername,
				"-p="+npmPassword,
				"-e="+npmEmail,
			)
			Expect(err).NotTo(HaveOccurred())
		}

		err = test.WasmeCliSplit("push " + imageName)
		Expect(err).NotTo(HaveOccurred())
	})
})
