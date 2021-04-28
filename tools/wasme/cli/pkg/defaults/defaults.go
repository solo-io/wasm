package defaults

import (
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/opts"
	"github.com/solo-io/wasm/tools/wasme/pkg/cache"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
)

func NewDefaultCacheWithAuth(opts *opts.AuthOptions) cache.Cache {
	// Pull command from a private registry still needs authorizer
	res, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
	puller := pull.NewPuller(res)

	return cache.NewCache(puller)
}
