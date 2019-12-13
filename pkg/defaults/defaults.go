package defaults

import (
	"github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
)

func NewDefaultCache() cache.Cache {
	resolver := resolver.NewResolver("", "", true, false)
	puller := pull.NewPuller(resolver)

	return &cache.CacheImpl{
		Puller:   puller,
		Resolver: resolver,
	}
}
