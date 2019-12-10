package cmd

import (
	"github.com/solo-io/wasme/pkg/cmd/build"
	"github.com/solo-io/wasme/pkg/cmd/initialize"
	"github.com/solo-io/wasme/pkg/cmd/list"
	"github.com/solo-io/wasme/pkg/version"
	"os"

	"github.com/solo-io/wasme/pkg/cmd/cache"
	"github.com/solo-io/wasme/pkg/cmd/catalog"
	"github.com/solo-io/wasme/pkg/cmd/login"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/cmd/pull"
	"github.com/solo-io/wasme/pkg/cmd/push"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "wasme [command]",
		Version: version.Version,
	}
	var opts opts.GeneralOptions
	cmd.AddCommand(
		initialize.InitCmd(),
		build.BuildCmd(),
		push.PushCmd(&opts),
		pull.PullCmd(&opts),
		cache.CacheCmd(&opts),
		catalog.CatalogCmd(&opts),
		login.LoginCmd(&opts),
		list.ListCmd(),
	)
	cmd.PersistentFlags().StringArrayVarP(&opts.Configs, "config", "c", nil, "auth config path")
	cmd.PersistentFlags().StringVarP(&opts.Username, "username", "u", "", "registry username")
	cmd.PersistentFlags().StringVarP(&opts.Password, "password", "p", "", "registry password")
	cmd.PersistentFlags().BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	cmd.PersistentFlags().BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "use plain http and not https")

	return cmd
}

func Run() {
	if err := Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
