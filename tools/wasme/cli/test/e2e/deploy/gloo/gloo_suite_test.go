package gloo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/solo-io/solo-kit/test/helpers"
)

func TestOperator(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	RunSpecs(t, "Deploy Gloo Suite")
}
