package build_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/build"
)

var _ = Describe("Env Proxy Args Passthrough", func() {
	It("should not pass through any extra env vars if none are set", func() {
		result := GetProxyEnvArgs()
		Expect(result).To(HaveLen(0), "shouldn't generate extra args")
	})

	It("should pass a single env var when set", func() {
		os.Setenv("http_proxy", "http://example.com")
		result := GetProxyEnvArgs()
		Expect(result).To(HaveLen(2))
		Expect(result[0]).To(Equal("-e"))
		Expect(result[1]).To(Equal("http_proxy=http://example.com"))
	})

	It("should pass multiple env vars when set", func() {
		os.Setenv("http_proxy", "http://example.com")
		os.Setenv("https_proxy", "https://example.com")
		os.Setenv("no_proxy", "https://example.com/foo")
		os.Setenv("GOPROXY", "https://example.com/bar")
		result := GetProxyEnvArgs()
		Expect(result).To(HaveLen(8))
		Expect(result[0]).To(Equal("-e"))
		Expect(result[1]).To(Equal("http_proxy=http://example.com"))
		Expect(result[2]).To(Equal("-e"))
		Expect(result[3]).To(Equal("https_proxy=https://example.com"))
		Expect(result[4]).To(Equal("-e"))
		Expect(result[5]).To(Equal("no_proxy=https://example.com/foo"))
		Expect(result[6]).To(Equal("-e"))
		Expect(result[7]).To(Equal("GOPROXY=https://example.com/bar"))
	})
})
