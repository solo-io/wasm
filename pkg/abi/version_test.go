package abi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/wasme/pkg/abi"
)

var _ = Describe("ABI Version Registry", func() {
	It("matches a real platform version with a registered platform version", func() {
		err := DefaultRegistry.ValidateIstioVersion(Version_6d525c67f39b36cdff9d688697f266c1b55e9cb7.Name, "1.4.2")
		Expect(err).NotTo(HaveOccurred())
		err = DefaultRegistry.ValidateIstioVersion(Version_6d525c67f39b36cdff9d688697f266c1b55e9cb7.Name, "1.4.32")
		Expect(err).NotTo(HaveOccurred())
		err = DefaultRegistry.ValidateIstioVersion("invalid_abiversion", "1.4.32")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("abi version invalid_abiversion not found"))
		err = DefaultRegistry.ValidateIstioVersion(Version_6d525c67f39b36cdff9d688697f266c1b55e9cb7.Name, "1.5.32")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no versions of istio found which match abi version"))
	})
})
