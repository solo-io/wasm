package defaults

import (
	"os"
	"path/filepath"

	"github.com/solo-io/wasme/pkg/cache"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/pull"
	"github.com/solo-io/wasme/pkg/resolver"
)

func NewDefaultCache(opts *opts.AuthOptions) cache.Cache {
	// Pull command from a private registry still needs authorizer
	resolver, _ := resolver.NewResolver(opts.Username, opts.Password, opts.Insecure, opts.PlainHTTP, opts.CredentialsFiles...)
	puller := pull.NewPuller(resolver)

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
