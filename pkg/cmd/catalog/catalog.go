package catalog

import (
	"context"

	"github.com/solo-io/extend-envoy/pkg/catalog"
	"github.com/solo-io/extend-envoy/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

type catalogOptions struct {
	*opts.GeneralOptions
}

func CatalogCmd(generalOptions *opts.GeneralOptions) *cobra.Command {
	var opts catalogOptions
	opts.GeneralOptions = generalOptions
	cmd := &cobra.Command{
		Use:   "catalog name[:tag|@digest] ...",
		Short: "interact with catalog",
		Long: `catalog
`,
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalog(opts)
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add name[:tag|@digest] ...",
		Short: "add to catalog",
		Long: `add
`,
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalog(opts)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "login name[:tag|@digest] ...",
		Short: "login to catalog",
		Long: `login
`,
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	})

	return cmd
}

func runCatalog(opts catalogOptions) error {

	return catalog.UpdateCatalogItem(context.Background(),
		"yuval123", "testrepo", "yuval.foo", "foo: bar")

}

func runLogin(opts catalogOptions) error {
	return catalog.Login(context.Background())
}
