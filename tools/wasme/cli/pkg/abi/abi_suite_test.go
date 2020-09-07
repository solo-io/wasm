package abi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAbi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Abi Suite")
}
