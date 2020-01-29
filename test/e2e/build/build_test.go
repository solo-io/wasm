package build_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/wasme/test"
)

var _ = Describe("Build", func() {
	It("builds and pushes the image", func() {
		err := test.RunMake("builder-image")
		Expect(err).NotTo(HaveOccurred())

		imageName := os.Getenv("FILTER_IMAGE_TAG")
		if imageName == "" {
			Skip("Skipping build/push test. To enable, set FILTER_IMAGE_TAG to the tag to use for the built/pushed image")
		}

		err = test.WasmeCliSplit("init test-filter --platform istio --platform-version 1.4.x --language cpp")
		Expect(err).NotTo(HaveOccurred())

		err = test.WasmeCliSplit("build -t " + imageName + " test-filter")
		Expect(err).NotTo(HaveOccurred())

		err = test.WasmeCliSplit("push " + imageName)
		Expect(err).NotTo(HaveOccurred())
	})
})
