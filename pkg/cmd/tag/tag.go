package tag

import (
	"context"

	"github.com/solo-io/wasme/pkg/model"
	"github.com/solo-io/wasme/pkg/store"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.StandardLogger()

type tagOptions struct {
	sourceImage string
	targetImage string

	storageDir string
}

func TagCmd(ctx *context.Context) *cobra.Command {
	var opts tagOptions
	cmd := &cobra.Command{
		Use:   "tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG]",
		Short: "Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.sourceImage = args[0]
			opts.targetImage = args[1]
			return runTag(*ctx, opts)
		},
	}

	cmd.Flags().StringVar(&opts.storageDir, "store", "", "Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store")
	return cmd
}

// override the tag of an image
type taggedImage struct {
	ref string
	model.Image
}

func (i *taggedImage) Ref() string {
	return i.ref
}

func runTag(ctx context.Context, opts tagOptions) error {
	imageStore := store.NewStore(opts.storageDir)

	sourceRef, err := model.FullRef(opts.sourceImage)
	if err != nil {
		return err
	}

	targetRef, err := model.FullRef(opts.targetImage)
	if err != nil {
		return err
	}

	image, err := imageStore.Get(sourceRef)
	if err != nil {
		return err
	}

	targetImage := &taggedImage{
		ref:   targetRef,
		Image: image,
	}

	if err := imageStore.Add(ctx, targetImage); err != nil {
		return err
	}

	// need to read filter to generate descriptor
	descriptor, err := targetImage.Descriptor()
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"digest": descriptor.Digest.String(),
		"image":  image.Ref(),
	}).Info("tagged image")

	return nil
}
