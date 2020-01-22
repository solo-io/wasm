package pull

import (
	"context"
	"io"

	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/model"
)

// represents the descriptors for an image, as well as accessors to the image contents
type Image interface {
	// ref to the image
	Ref() string

	// get the image (Filter) descriptor
	Descriptor() (ocispec.Descriptor, error)

	// get the filter .wasm file from the image
	FetchFilter(ctx context.Context) (io.ReadCloser, error)

	// get the filter config from the image
	FetchConfig(ctx context.Context) (*config.Config, error)
}

type imageDescriptors struct {
	children []ocispec.Descriptor
	ref      string
	resolver remotes.Resolver
}

func (i *imageDescriptors) Ref() string {
	return i.ref
}

func (i *imageDescriptors) Descriptor() (ocispec.Descriptor, error) {
	return i.getDescriptor(model.CodeMediaType)
}

func (i *imageDescriptors) FetchFilter(ctx context.Context) (io.ReadCloser, error) {
	desc, err := i.getDescriptor(model.CodeMediaType)
	if err != nil {
		return nil, err
	}

	fetcher, err := i.resolver.Fetcher(ctx, i.ref)
	if err != nil {
		return nil, err
	}

	return fetcher.Fetch(ctx, desc)
}

func (i *imageDescriptors) FetchConfig(ctx context.Context) (*config.Config, error) {
	desc, err := i.getDescriptor(model.ConfigMediaType)
	if err != nil {
		return nil, err
	}

	fetcher, err := i.resolver.Fetcher(ctx, i.ref)
	if err != nil {
		return nil, err
	}

	rc, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}

	return config.FromReader(rc)
}

func (i *imageDescriptors) getDescriptor(mediaType string) (ocispec.Descriptor, error) {
	for _, child := range i.children {
		if child.MediaType == mediaType {
			return child, nil
		}
	}
	return ocispec.Descriptor{}, errors.Errorf("media type %v not found on image", mediaType)
}
