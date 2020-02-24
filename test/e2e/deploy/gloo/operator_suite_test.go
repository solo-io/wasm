package gloo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGlooDeploy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deploy Gloo Suite")
}
