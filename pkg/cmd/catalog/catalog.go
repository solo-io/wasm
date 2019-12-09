package catalog

import (
	"context"

	"github.com/solo-io/wasme/pkg/auth"
	"github.com/solo-io/wasme/pkg/catalog"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

type catalogOptions struct {
	*opts.GeneralOptions
}

func CatalogCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts catalogOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "interact with catalog",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add name[:tag|@digest] ...",
		Short: "add to catalog",
		Long: `add
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalog(opts, args[0])
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "login to catalog",
		Long: `login allows you pushing to getwasm.io and automate the process of 
		creating a PR to publish you content to the hub.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	})

	return cmd
}

func runCatalog(opts catalogOptions, ref string) error {

	return catalog.UpdateCatalogItem(context.Background(), ref)

}

func runLogin(opts catalogOptions) error {
	return auth.Login(context.Background())
}
