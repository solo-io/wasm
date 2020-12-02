package build

import (
	"fmt"
	"os"
)

var passThroughVars = []string{
	"GOPROXY",
	"http_proxy",
	"https_proxy",
	"no_proxy",
}

// getProxyEnvArgs reads several environment variables and returns
// the arguments to pass them into the docker container used
// during a wasme build command
func getProxyEnvArgs() []string {
	var proxyEnvArgs []string
	for _, envVar := range passThroughVars {
		val, isSet := os.LookupEnv(envVar)
		if isSet {
			proxyEnvArgs = append(proxyEnvArgs, "-e", fmt.Sprintf("%s=%s", envVar, val))
		}
	}
	return proxyEnvArgs
}
