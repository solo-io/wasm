package build

import (
	"context"

	"github.com/spf13/cobra"
)

func rustCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	var bazel bazelOptions
	cmd := &cobra.Command{
		Use:   "rust SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag>",
		Short: "Build a wasm image from a Rust filter using Bazel-in-Docker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceDir = args[0]
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return runBazelBuild(*build, bazel)
			})
		},
	}

	cmd.Flags().StringVarP(&bazel.buildDir, "build-dir", "b", ".", "Directory containing the target BUILD file.")
	cmd.Flags().StringVarP(&bazel.bazelOutput, "bazel-output", "f", "filter.wasm", "Path relative to `bazel-bin` to the wasm file produced by running the Bazel target.")
	cmd.Flags().StringVarP(&bazel.bazelTarget, "bazel-target", "g", ":filter", "Name of the bazel target to run.")
	return cmd
}
