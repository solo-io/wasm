package gloo_test

import (
	"testing"

	"github.com/solo-io/solo-kit/test/helpers"

	. "github.com/onsi/ginkgo"
)

func TestOperator(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	RunSpecs(t, "Deploy Gloo Suite")
}
