package cmd

import (
	"context"
	"os"

	"github.com/solo-io/wasm/tools/wasme/pkg/defaults"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/tag"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/operator"

	"github.com/sirupsen/logrus"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/build"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/deploy"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/initialize"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/list"
	"github.com/solo-io/wasm/tools/wasme/pkg/version"

	ctxo "github.com/deislabs/oras/pkg/context"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/cache"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/login"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/opts"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/pull"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/cmd/push"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var auth opts.AuthOptions
	var general opts.GeneralOptions

	ctx2 := context.Background()
	ctx := &ctx2
	cmd := &cobra.Command{
		Use:     "wasme [command]",
		Short:   "The tool for building, pushing, and deploying Envoy WebAssembly Filters",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if general.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				ctx2 := ctxo.WithLoggerDiscarded(*ctx)
				*ctx = ctx2
			}
			// set default auth configs
			if len(auth.CredentialsFiles) == 0 {
				auth.CredentialsFiles = []string{defaults.WasmeCredentialsFile}
			}
		},
	}

	commandsWithAuth := []*cobra.Command{
		push.PushCmd(ctx, &auth),
		pull.PullCmd(ctx, &auth),
		cache.CacheCmd(ctx, &auth),
	}

	for _, cmd := range commandsWithAuth {
		auth.AddToFlags(cmd.PersistentFlags())
	}

	commands := append(commandsWithAuth,
		initialize.InitCmd(),
		build.BuildCmd(ctx),
		login.LoginCmd(),
		list.ListCmd(),
		deploy.DeployCmd(ctx, cmd.PersistentPreRun),
		deploy.UndeployCmd(ctx),
		operator.OperatorCmd(ctx),
		tag.TagCmd(ctx))

	cmd.AddCommand(
		commands...,
	)

	general.AddToFlags(cmd.PersistentFlags())

	return cmd
}

func Run() {
	if err := Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
