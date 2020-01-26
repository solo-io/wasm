package store

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/model"
)

type Image interface {
	model.Image
}

// an image stored on disk
type storedImage struct {
	ref        string
	descriptor ocispec.Descriptor
	filter     func() (model.Filter, error)
	config     *config.Config
}

func NewStorableImage(ref string, descriptor ocispec.Descriptor, filter func() (model.Filter, error), config *config.Config) *storedImage {
	return &storedImage{ref: ref, descriptor: descriptor, filter: filter, config: config}
}

func (i *storedImage) Ref() string {
	return i.ref
}

func (i *storedImage) Descriptor() (ocispec.Descriptor, error) {
	return i.descriptor, nil
}

func (i *storedImage) FetchFilter(ctx context.Context) (model.Filter, error) {
	return i.filter()
}

func (i *storedImage) FetchConfig(ctx context.Context) (*config.Config, error) {
	return i.config, nil
}
