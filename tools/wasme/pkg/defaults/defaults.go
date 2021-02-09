package defaults

import (
	"context"

	"github.com/solo-io/wasm/tools/wasme/pkg/cache"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
)

// Trigger CI / Don't merge this.

func NewDefaultCache() cache.Cache {
	return cache.NewCache(NewDefaultPuller())
}

func NewDefaultCacheWithContext(ctx context.Context) cache.Cache {
	return cache.NewCacheWithConext(ctx, NewDefaultPuller())
}

func NewDefaultPuller() pull.ImagePuller {
	// Can pull from non-private registries
	res, _ := resolver.NewResolver("", "", true, false)
	return pull.NewPuller(res)
}
