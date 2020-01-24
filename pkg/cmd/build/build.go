package build

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/model"
	"github.com/solo-io/wasme/pkg/store"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/version"

	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type buildOptions struct {
	sourceDir    string
	configFile   string
	tag          string
	storageDir   string
	builderImage string
	buildDir     string
	bazelOutput  string
	bazelTarget  string
}

func BuildCmd() *cobra.Command {
	var opts buildOptions
	cmd := &cobra.Command{
		Use:   "build SOURCE_DIRECTORY [-b <bazel target>] [-t <name:tag>]",
		Short: "Build a wasm image from the filter source directory using Bazel-in-Docker",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("invalid number of arguments")
			}
			opts.sourceDir = args[0]
			return runBuild(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "The image ref with which to tag this image. Specified in the format <name:tag>. Required")
	cmd.Flags().StringVarP(&opts.configFile, "config", "c", "", "The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/filter-config.json. This file must be present in order to build the image.")
	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")
	cmd.Flags().StringVarP(&opts.builderImage, "image", "i", "quay.io/solo-io/ee-builder:"+version.Version, "Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image")
	cmd.Flags().StringVarP(&opts.buildDir, "build-dir", "b", ".", "Directory containing the target BUILD file.")
	cmd.Flags().StringVarP(&opts.bazelOutput, "bazel-ouptut", "f", "filter.wasm", "Path relative to `bazel-bin` to the wasm file produced by running the Bazel target.")
	cmd.Flags().StringVarP(&opts.bazelTarget, "bazel-target", "t", ":filter.wasm", "Name of the bazel target to run.")
	return cmd
}

func runBuild(opts buildOptions) error {
	configFile := opts.configFile
	if configFile == "" {
		configFile = filepath.Join(opts.sourceDir, "filter-config.json")
	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	cfg, err := config.FromBytes(configBytes)
	if err != nil {
		return err
	}

	sourceDir, err := filepath.Abs(opts.sourceDir)
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("/tmp", "wasme")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	// container paths are currently hard-coded in builder image
	args := []string{
		"run",
		"--rm",
		"-v", sourceDir + ":/src/workspace",
		"-v", tmpDir + ":/tmp/build_output",
		"-w", "/src/workspace",
		"-e", "BUILD_BASE=" + opts.buildDir,
		"-e", "BAZEL_OUTPUT=" + opts.bazelOutput,
		"-e", "TARGET=" + opts.bazelTarget,
		opts.builderImage,
	}

	log.WithFields(logrus.Fields{
		"args": args,
	}).Info("running bazel-in-docker build...")

	if err := docker(os.Stdout, os.Stderr, args...); err != nil {
		return err
	}

	// filter.wasm currently hard-coded in bazel BUILD file
	tmpFile := filepath.Join(tmpDir, "filter.wasm")

	log.WithFields(logrus.Fields{
		"tmp_file": tmpFile,
		"tag":      opts.tag,
	}).Info("adding image to cache...")

	filterFile, err := os.Open(tmpFile)
	if err != nil {
		return err
	}

	descriptor, err := model.GetDescriptor(filterFile)
	if err != nil {
		return err
	}

	// expand the ref to contain :latest suffix if no tag provided
	imageRef := func() string {
		parts := strings.Split(opts.tag, ":")
		if len(parts) == 2 {
			return opts.tag
		}
		return opts.tag + ":latest"
	}()

	image := store.NewStorableImage(imageRef, descriptor, filterFile, cfg)

	if err := store.NewStore(opts.storageDir).Add(context.Background(), image); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"digest": descriptor.Digest.String(),
		"image":  imageRef,
	}).Info("tagged image")

	return nil
}

func docker(stdout, stderr io.Writer, args ...string) error {
	return execCmd(stdout, stderr, "docker", args...)
}

func execCmd(stdout, stderr io.Writer, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stderr = stderr
	command.Stdout = stdout
	return command.Run()
}
