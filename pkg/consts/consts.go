package consts

import "os"

var (
	HubDomain = func() string {
		if customDomain := os.Getenv("IMAGE_REGISTRY"); customDomain != "" {
			return customDomain
		}
		return "localhost:8080"
		//return "webassemblyhub.io"
	}()
)
