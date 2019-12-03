package catalog

import (
	"context"
	"os"

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
		Short: "submit to catalog",
		Long: `catalog
`,
		//Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCatalog(opts)
		},
	}

	return cmd
}

func runCatalog(opts catalogOptions) error {

	return catalog.UpdateCatalogItem(context.Background(), os.Getenv("GITHUB_API_TOKEN"),
		"yuval123", "testrepo", "yuval.foo", "foo: bar")

}
