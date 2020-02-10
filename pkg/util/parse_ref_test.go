package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/wasme/pkg/util"
)

var _ = Describe("ParseRef", func() {
	It("parses an image ref", func() {
		ref := "localhost:8080/taco/tuesdays:v1"
		name, tag, err := SplitImageRef(ref)
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal("localhost:8080/taco/tuesdays"))
		Expect(tag).To(Equal("v1"))
	})
	It("parses an image ref", func() {
		ref := "localhost:8080/taco/tuesdays"
		name, tag, err := SplitImageRef(ref)
		Expect(err).NotTo(HaveOccurred())
		Expect(name).To(Equal("localhost:8080/taco/tuesdays"))
		Expect(tag).To(Equal("latest"))
	})
})
