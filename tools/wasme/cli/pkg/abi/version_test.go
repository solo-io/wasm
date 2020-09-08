package abi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/wasm/tools/wasme/cli/pkg/abi"
)

var _ = Describe("ABI Version Registry", func() {
	It("matches a real platform version with a registered platform version", func() {
		err := DefaultRegistry.ValidateIstioVersion([]string{Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f.Name}, "1.5.0")
		Expect(err).NotTo(HaveOccurred())
		err = DefaultRegistry.ValidateIstioVersion([]string{Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f.Name}, "1.5.32")
		Expect(err).NotTo(HaveOccurred())
		err = DefaultRegistry.ValidateIstioVersion([]string{"invalid_abiversion"}, "1.5.32")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("abi versions [invalid_abiversion] not found"))
		err = DefaultRegistry.ValidateIstioVersion([]string{Version_edc016b1fa5adca3ebd3d7020eaed0ad7b8814ca.Name}, "1.5.32")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("no versions of istio found which support abi versions"))
		err = DefaultRegistry.ValidateIstioVersion([]string{Version_097b7f2e4cc1fb490cc1943d0d633655ac3c522f.Name}, "1.6.0")
		Expect(err).NotTo(HaveOccurred())
	})
})
