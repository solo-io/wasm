package build

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/solo-io/wasm/tools/wasme/pkg/config"
	"github.com/solo-io/wasm/tools/wasme/pkg/model"
	"github.com/solo-io/wasm/tools/wasme/pkg/store"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasm/tools/wasme/cli/pkg/version"

	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type buildOptions struct {
	sourceDir    string
	configFile   string
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
		Long:  `Options for the build are specific to the target language.`,
	}

	cmd.PersistentFlags().StringVarP(&opts.tag, "tag", "t", "", "The image ref with which to tag this image. Specified in the format <name:tag>. Required")
	cmd.PersistentFlags().StringVarP(&opts.configFile, "config", "c", "", "The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.")
	cmd.PersistentFlags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")
	cmd.PersistentFlags().StringVarP(&opts.builderImage, "image", "i", "quay.io/solo-io/ee-builder:"+version.Version, "Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image")
	cmd.PersistentFlags().StringVarP(&opts.tmpDir, "tmp-dir", "", "", "Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux. If unset, temporary files will be removed after build")

	cmd.AddCommand(
		cppCmd(ctx, &opts),
		assemblyscriptCmd(ctx, &opts),
		precompiledCmd(ctx, &opts),
	)

	return cmd
}

func runBuild(ctx context.Context, opts *buildOptions, getFilter func(opts *buildOptions) (string, error)) error {
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

	tmpDir := opts.tmpDir
	customTmpDir := true
	if tmpDir != "" {
		customTmpDir = false
	} else {
		// workaround for darwin, cannot mount /var to docker
		if runtime.GOOS == "darwin" {
			tmpDir = "/tmp"
		}
	}

	if customTmpDir {
		tmpDir, err = ioutil.TempDir(tmpDir, "wasme")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)
	}

	// use abs dir because docker requires it
	tmpDir, err = filepath.Abs(tmpDir)
	if err != nil {
		return err
	}
	opts.tmpDir = tmpDir

	filterFile, err := getFilter(opts)
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

func getDescriptor(filterBytes []byte) (ocispec.Descriptor, error) {
	descriptor, err := model.GetDescriptor(bytes.NewBuffer(filterBytes))
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return descriptor, nil
}
