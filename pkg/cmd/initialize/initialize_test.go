package initialize

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Initialize", func() {
	It("works", func() {
		err := runInit(initOptions{destDir:"./foo"})
		Expect(err).NotTo(HaveOccurred())
	})
})
