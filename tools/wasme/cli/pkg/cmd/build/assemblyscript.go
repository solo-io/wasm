package build

import (
	"context"
	"os"
	"path/filepath"

	"github.com/solo-io/wasm/tools/wasme/cli/pkg/defaults"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type npmOpts struct {
	username, password, email string
}

func assemblyscriptCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	var npm npmOpts
	cmd := &cobra.Command{
		Use:   "assemblyscript SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag>",
		Short: "Build a wasm image from an AssemblyScript filter using NPM-in-Docker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceDir = args[0]
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return runNpmBuild(*build, npm)
			})
		},
	}

	cmd.PersistentFlags().StringVarP(&npm.username, "username", "u", "", "Username for logging in to NPM before running npm install. Optional")
	cmd.PersistentFlags().StringVarP(&npm.password, "password", "p", "", "Password for logging in to NPM before running npm install. Optional")
	cmd.PersistentFlags().StringVarP(&npm.email, "email", "e", "", "Email for logging in to NPM before running npm install. Optional")

	return cmd
}

func runNpmBuild(build buildOptions, npm npmOpts) (string, error) {
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
		"-e", "BUILD_TOOL=npm", // required by build-filter.sh in container
	}

	if npm.username != "" && npm.password != "" && npm.email != "" {
		args = append(args, "-e", "NPM_USERNAME="+npm.username, "-e", "NPM_PASSWORD="+npm.password, "-e", "NPM_EMAIL="+npm.email)
	}

	args = append(args, defaults.GetProxyEnvArgs()...)

	log.WithFields(logrus.Fields{
		"args": args,
	}).Debug("running npm-in-docker build...")

	if err := util.DockerRun(os.Stdout, os.Stderr, nil, build.builderImage, args, nil); err != nil {
		return "", err
	}

	// filter.wasm currently hard-coded in package.json file
	return filepath.Join(build.tmpDir, "filter.wasm"), nil
}
