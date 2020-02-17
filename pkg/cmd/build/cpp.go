package build

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type bazelOptions struct {
	buildDir    string
	bazelOutput string
	bazelTarget string
}

func cppCmd(ctx *context.Context, opts *buildOptions) *cobra.Command {
	var bazel bazelOptions
	cmd := &cobra.Command{
		Use:   "cpp SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag>",
		Short: "Build a wasm image from a CPP filter using Bazel-in-Docker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceDir = args[0]
			return runBuild(*ctx, opts, func(build *buildOptions) (s string, err error) {
				return runBazelBuild(*build, bazel)
			})
		},
	}

	cmd.Flags().StringVarP(&bazel.buildDir, "build-dir", "b", ".", "Directory containing the target BUILD file.")
	cmd.Flags().StringVarP(&bazel.bazelOutput, "bazel-ouptut", "f", "filter.wasm", "Path relative to `bazel-bin` to the wasm file produced by running the Bazel target.")
	cmd.Flags().StringVarP(&bazel.bazelTarget, "bazel-target", "g", ":filter.wasm", "Name of the bazel target to run.")
	return cmd
}

func runBazelBuild(build buildOptions, bazel bazelOptions) (string, error) {
	sourceDir, err := filepath.Abs(build.sourceDir)
	if err != nil {
		return "", err
	}

	tmpDirName := build.tmpDir
	// workaround for darwin, cannot mount /var to docker
	if tmpDirName == "" && runtime.GOOS == "darwin" {
		tmpDirName = "/tmp"
	}
	tmpDir, err := ioutil.TempDir(tmpDirName, "wasme")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// use abs dir because docker requires it
	tmpDir, err = filepath.Abs(tmpDir)
	if err != nil {
		return "", err
	}

	// container paths are currently hard-coded in builder image
	args := []string{
		"run",
		"--rm",
		"-v", sourceDir + ":/src/workspace",
		"-v", tmpDir + ":/build_output",
		"-w", "/src/workspace",
		"-e", "BUILD_BASE=" + bazel.buildDir,
		"-e", "BAZEL_OUTPUT=" + bazel.bazelOutput,
		"-e", "TARGET=" + bazel.bazelTarget,
		"-e", "BUILD_TOOL=bazel", // required by build-filter.sh in container
		build.builderImage,
	}

	log.WithFields(logrus.Fields{
		"args": args,
	}).Debug("running npm-in-docker build...")

	if err := docker(os.Stdout, os.Stderr, args...); err != nil {
		return "", err
	}

	// filter.wasm currently hard-coded in bazel BUILD file
	return filepath.Join(tmpDir, "filter.wasm"), nil
}
