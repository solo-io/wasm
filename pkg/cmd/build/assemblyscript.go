package build

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func assemblyscriptCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assemblyscript SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag>",
		Short: "Build a wasm image from an AssemblyScript filter using NPM-in-Docker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceDir = args[0]
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return runNpmBuild(*build)
			})
		},
	}

	return cmd
}

func runNpmBuild(build buildOptions) (string, error) {
	sourceDir, err := filepath.Abs(build.sourceDir)
	if err != nil {
		return "", err
	}

	// container paths are currently hard-coded in builder image
	args := []string{
		"run",
		"--rm",
		"-v", sourceDir + ":/src/workspace",
		"-v", build.tmpDir + ":/build_output",
		"-w", "/src/workspace",
		"-e", "BUILD_TOOL=npm", // required by build-filter.sh in container
		build.builderImage,
	}

	log.WithFields(logrus.Fields{
		"args": args,
	}).Info("running npm-in-docker build...")

	if err := docker(os.Stdout, os.Stderr, args...); err != nil {
		return "", err
	}

	// filter.wasm currently hard-coded in bazel BUILD file
	return filepath.Join(build.tmpDir, "filter.wasm"), nil
}
