package build_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/log"
)

func TestBuild(t *testing.T) {
	if os.Getenv("FILTER_BUILD_IMAGE_TAG") == "" {
		log.Warnf("This test is disabled. " +
			"To enable, set FILTER_BUILD_IMAGE_TAG in your env.")
		return
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Suite")
}
