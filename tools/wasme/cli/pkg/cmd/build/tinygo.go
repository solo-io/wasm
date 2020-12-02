package build

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"
	"github.com/spf13/cobra"
)

func tinygoCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tinygo SOURCE_DIRECTORY -t <name:tag>",
		Short: "Build a wasm image from a TinyGo filter using TinyGo-in-Docker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceDir = args[0]
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return runTinyGoBuild(*build)
			})
		},
	}
	return cmd
}

func runTinyGoBuild(build buildOptions) (string, error) {
	sourceDir, err := filepath.Abs(build.sourceDir)
	if err != nil {
		return "", err
	}

	// container paths are currently hard-coded in builder image
	args := []string{
		"--rm",
		"-v", sourceDir + ":/src/workspace",
		"-v", build.tmpDir + ":/build_output",
		"-w", "/src/workspace",
		"-e", "BUILD_TOOL=tinygo", // required by build-filter.sh in container
	}

	args = append(args, GetProxyEnvArgs()...)

	log.WithFields(logrus.Fields{
		"args": args,
	}).Debug("running TinyGo-in-docker build...")

	if err := util.DockerRun(os.Stdout, os.Stderr, nil, build.builderImage, args, nil); err != nil {
		return "", err
	}

	// filter.wasm currently hard-coded in package.json file
	return filepath.Join(build.tmpDir, "filter.wasm"), nil
}
