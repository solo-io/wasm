package consts

import "os"

var (
	HubDomain = func() string {
		if customDomain := os.Getenv("WASM_IMAGE_REGISTRY"); customDomain != "" {
			return customDomain
		}
		return "webassemblyhub.io"
	}()
)

const (
	CachePort = 9979
)
