package login

import (
	"context"

	"github.com/solo-io/wasme/pkg/auth"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

type catalogOptions struct {
	*opts.GeneralOptions
}

func LoginCmd(ctx *context.Context, generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts catalogOptions
	opts.GeneralOptions = generalOptions
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

	return cmd
}

func runLogin(ctx context.Context, opts catalogOptions) error {
	return auth.Login(ctx)
}
