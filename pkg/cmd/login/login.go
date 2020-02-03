package login

import (
	"context"

	"github.com/solo-io/wasme/pkg/auth"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	skipBrowser bool
}

func LoginCmd(ctx *context.Context) *cobra.Command {
	var opts loginOptions
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login so you can push images to webassemblyhub.io and submit them to the curated catalog",
		Long: `login allows you pushing to webassemblyhub.io and automate the process of 
		creating a PR to publish you content to the hub.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(*ctx, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.skipBrowser, "skip-browser", false, "skip opening the default browser to the "+
		"login URL")

	return cmd
}

func runLogin(ctx context.Context, opts loginOptions) error {
	return auth.Login(ctx, opts.skipBrowser)
}
