package build

import (
	"context"

	"github.com/spf13/cobra"
)

func precompiledCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "precompiled COMPILED_FILTER_FILE --tag <name:tag> --config <image config>",
		Short: "Build a wasm image from a Precompiled filter.",
		Long: `
wasme supports building deployable images from a precompiled .wasm file. The user must provide their own configuration file with the --config flag.

The specification for this flag can be found here: [{{< versioned_link_path fromRoot="/reference/image_config">}}]({{< versioned_link_path fromRoot="/reference/image_config">}})
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return args[0], nil
			})
		},
	}

	return cmd
}
