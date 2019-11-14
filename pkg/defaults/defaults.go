package defaults

import (
	"github.com/solo-io/extend-envoy/pkg/cache"
	"github.com/solo-io/extend-envoy/pkg/pull"
	"github.com/solo-io/extend-envoy/pkg/resolver"
)

func NewDefaultCache() cache.Cache {
	resolver := resolver.NewResolver("", "", true, false)
	puller := pull.NewPuller(resolver)

	return &cache.CacheImpl{
		Puller:   puller,
		Resolver: resolver,
	}
}
