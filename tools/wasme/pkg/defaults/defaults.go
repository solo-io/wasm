package defaults

import (
	"github.com/solo-io/wasm/tools/wasme/pkg/cache"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
)

func NewDefaultCache() cache.Cache {
	// Can pull from non-private registries
	res, _ := resolver.NewResolver("", "", true, false)
	puller := pull.NewPuller(res)

	return cache.NewCache(puller)
}
