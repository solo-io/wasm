package defaults

import (
	"os"
	"path/filepath"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/opts"
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

func NewDefaultCacheWithAuth(opts *opts.AuthOptions) cache.Cache {
	// Pull command from a private registry still needs authorizer
	res, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
	puller := pull.NewPuller(res)

	return cache.NewCache(puller)
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
