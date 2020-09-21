package gloo_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
)

func TestOperator(t *testing.T) {
	if os.Getenv("FILTER_IMAGE_GLOO_TAG") == "" {
		log.Warnf("This test is disabled. " +
			"To enable, set FILTER_IMAGE_GLOO_TAG to a valid gloo wasm filter image in your env.")
		return
	}
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	RunSpecs(t, "Deploy Gloo Suite")
}
