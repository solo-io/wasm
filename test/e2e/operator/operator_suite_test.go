package operator_test

import (
	"testing"

	"github.com/solo-io/solo-kit/test/helpers"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

func TestOperator(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	helpers.SetupLog()
	RunSpecs(t, "Operator Suite")
}
