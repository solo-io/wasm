package defaults

import (
	"os"
	"path/filepath"

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

var (
	WasmeConfigDir       = os.Getenv("HOME") + "/.wasme"
	WasmeImageDir        = filepath.Join(WasmeConfigDir, "store")
	WasmeCredentialsFile = filepath.Join(WasmeConfigDir, "credentials.json")
)
