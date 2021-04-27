package defaults

import (
	"context"
	"os"
	"path/filepath"

	"github.com/solo-io/wasm/tools/wasme/pkg/cache"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"github.com/solo-io/wasm/tools/wasme/pkg/resolver"
)

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

var (
	WasmeConfigDir       = home() + "/.wasme"
	WasmeImageDir        = filepath.Join(WasmeConfigDir, "store")
	WasmeCredentialsFile = filepath.Join(WasmeConfigDir, "credentials.json")
)

func home() string {
	dir, _ := os.UserHomeDir()
	return dir
}
