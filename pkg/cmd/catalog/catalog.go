package catalog

import (
	"context"

	"github.com/solo-io/wasme/pkg/catalog"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

type catalogOptions struct {
	*opts.GeneralOptions
}

func CatalogCmd(ctx *context.Context, generalOptions *opts.GeneralOptions) *cobra.Command {
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
			return runCatalog(*ctx, opts, args[0])
		},
	})

	return cmd
}

func runCatalog(ctx context.Context, opts catalogOptions, ref string) error {

	return catalog.UpdateCatalogItem(ctx, ref)

}
