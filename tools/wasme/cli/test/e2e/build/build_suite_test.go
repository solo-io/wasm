package build_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/wasm/tools/wasme/cli/test"
)

func TestBuild(t *testing.T) {
	if test.GetBuildImageTag() == "" {
		log.Warnf("This test is disabled. " +
			"To enable, set FILTER_BUILD_IMAGE_TAG in your env.")
		return
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Suite")
}
