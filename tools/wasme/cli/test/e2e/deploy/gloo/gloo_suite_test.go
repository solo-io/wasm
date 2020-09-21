package gloo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/wasm/tools/wasme/cli/test"
)

func TestOperator(t *testing.T) {
	if test.GetImageTagGloo() == "" {
		log.Warnf("This test is disabled. " +
			"To enable, set FILTER_IMAGE_GLOO_TAG to a valid gloo wasm filter image in your env.")
		return
	}
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	RunSpecs(t, "Deploy Gloo Suite")
}
