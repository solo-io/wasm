package defaults

import (
	"github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
)

func NewDefaultCache() cache.Cache {
	// cache doesn't need authorizer as it doesn't push
	resolver, _ := resolver.NewResolver("", "", true, false)
	puller := pull.NewPuller(resolver)

	return cache.NewCache(puller)
}
