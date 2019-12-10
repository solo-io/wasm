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

func LoginCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts catalogOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login so you can push images to webassemblyhub.io and submit them to the curated catalog",
		Long: `login allows you pushing to webassemblyhub.io and automate the process of 
		creating a PR to publish you content to the hub.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	}

	return cmd
}

func runLogin(opts catalogOptions) error {
	return auth.Login(context.Background())
}
