package build

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

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
	wasmFile     string
	tag          string
	storageDir   string
	builderImage string
	tmpDir       string
}

func BuildCmd(ctx *context.Context) *cobra.Command {
	var opts buildOptions
	cmd := &cobra.Command{
		Use:   "build LANGUAGE SOURCE_DIRECTORY  -t <name:tag> [--options...]",
		Short: "Build a wasm image from the filter source directory.",
		Long: `Options for the build are specific to the target language.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if opts.wasmFile == "" {
					return fmt.Errorf("must provide either SOURCE_DIRECTORY or --wasm-file to build an image")
				}
			} else {
				opts.sourceDir = args[0]
			}
			return runBazelBuild(*ctx, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "The image ref with which to tag this image. Specified in the format <name:tag>. Required")
	cmd.Flags().StringVarP(&opts.configFile, "config", "c", "", "The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.")
	cmd.Flags().StringVarP(&opts.wasmFile, "wasm-file", "", "", "If specified, wasme will use the provided path to a compiled filter wasm to produce the image. The bazel build will be skipped and the wasm-file will be used instead.")
	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")
	cmd.Flags().StringVarP(&opts.builderImage, "image", "i", "quay.io/solo-io/ee-builder:"+version.Version, "Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image")
	cmd.Flags().StringVarP(&opts.tmpDir, "tmp-dir", "", "", "Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux.")
	return cmd
}

func runBuild(ctx context.Context, opts buildOptions, getFilter func() (string, error)) error {
	configFile := opts.configFile
	if configFile == "" {
		configFile = filepath.Join(opts.sourceDir, "runtime-config.json")
	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	cfg, err := config.FromBytes(configBytes)
	if err != nil {
		return err
	}

	filterFile, err := getFilter()
	if err != nil {
		return errors.Wrap(err, "failed producing filter file")
	}

	log.WithFields(logrus.Fields{
		"filter file": filterFile,
		"tag":         opts.tag,
	}).Info("adding image to cache...")

	filterBytes, err := ioutil.ReadFile(filterFile)
	if err != nil {
		return err
	}

	// need to read filter to generate descriptor
	descriptor, err := getDescriptor(filterBytes)
	if err != nil {
		return err
	}

	image, err := store.NewStorableImage(opts.tag, descriptor, filterBytes, cfg)
	if err != nil {
		return err
	}

	if err := store.NewStore(opts.storageDir).Add(ctx, image); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"digest": descriptor.Digest.String(),
		"image":  image.Ref(),
	}).Info("tagged image")

	return nil
}

type npmOpts struct {

}

func runNpmBuild(build buildOptions, bazel npmOpts) (string, error) {
		sourceDir, err := filepath.Abs(build.sourceDir)
		if err != nil {
			return "",  err
		}

		tmpDirName := build.tmpDir
		// workaround for darwin, cannot mount /var to docker
		if tmpDirName == "" && runtime.GOOS == "darwin" {
			tmpDirName = "/tmp"
		}
		tmpDir, err := ioutil.TempDir(tmpDirName, "wasme")
		if err != nil {
			return "",  err
		}
		defer os.RemoveAll(tmpDir)

		// use abs dir because docker requires it
		tmpDir, err = filepath.Abs(tmpDir)
		if err != nil {
			return "",  err
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
		}).Info("running npm-in-docker build...")

		if err := docker(os.Stdout, os.Stderr, args...); err != nil {
			return "",  err
		}

		// filter.wasm currently hard-coded in bazel BUILD file
		return filepath.Join(tmpDir, "filter.wasm"), nil
}

func runBazelBuild(ctx context.Context, build buildOptions, bazel bazelOptions) (string, error) {
		sourceDir, err := filepath.Abs(build.sourceDir)
		if err != nil {
			return "",  err
		}

		tmpDirName := build.tmpDir
		// workaround for darwin, cannot mount /var to docker
		if tmpDirName == "" && runtime.GOOS == "darwin" {
			tmpDirName = "/tmp"
		}
		tmpDir, err := ioutil.TempDir(tmpDirName, "wasme")
		if err != nil {
			return "",  err
		}
		defer os.RemoveAll(tmpDir)

		// use abs dir because docker requires it
		tmpDir, err = filepath.Abs(tmpDir)
		if err != nil {
			return "",  err
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
		}).Info("running npm-in-docker build...")

		if err := docker(os.Stdout, os.Stderr, args...); err != nil {
			return "",  err
		}

		// filter.wasm currently hard-coded in bazel BUILD file
		return filepath.Join(tmpDir, "filter.wasm"), nil
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

func getDescriptor(filterBytes []byte) (ocispec.Descriptor, error) {
	descriptor, err := model.GetDescriptor(bytes.NewBuffer(filterBytes))
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return descriptor, nil
}

type bazelOptions struct {
	buildDir     string
	bazelOutput  string
	bazelTarget  string
}

func bazelCmd(ctx *context.Context, opts buildOptions) *cobra.Command {

	var bazel bazelOptions
	cmd := &cobra.Command{
		Use:   "build SOURCE_DIRECTORY [-b <bazel target>] [-t <name:tag>]",
		Short: "Build a wasm image from the filter source directory using Bazel-in-Docker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if opts.wasmFile == "" {
					return fmt.Errorf("must provide either SOURCE_DIRECTORY or --wasm-file to build an image")
				}
			} else {
				opts.sourceDir = args[0]
			}
			return runBazelBuild(*ctx, opts)
		},
	}

	cmd.Flags().StringVarP(&bazel.buildDir, "build-dir", "b", ".", "Directory containing the target BUILD file.")
	cmd.Flags().StringVarP(&bazel.bazelOutput, "bazel-ouptut", "f", "filter.wasm", "Path relative to `bazel-bin` to the wasm file produced by running the Bazel target.")
	cmd.Flags().StringVarP(&bazel.bazelTarget, "bazel-target", "g", ":filter.wasm", "Name of the bazel target to run.")
	return cmd
}