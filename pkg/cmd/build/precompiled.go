package build

import (
	"context"
	"github.com/spf13/cobra"
)

func precompiledCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "precompiled COMPILED_FILTER_FILE -t <name:tag>",
		Short: "Build a wasm image from a Precompiled filter.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return args[0], nil
			})
		},
	}

	return cmd
}
