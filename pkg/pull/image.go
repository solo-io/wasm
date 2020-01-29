package pull

import (
	"context"

	"github.com/solo-io/wasme/pkg/model"

	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/config"
)

// an image that was pulled from a remote registry
type pulledImage struct {
	children []ocispec.Descriptor
	ref      string
	resolver remotes.Resolver
}

func (i *pulledImage) Ref() string {
	return i.ref
}

func (i *pulledImage) Descriptor() (ocispec.Descriptor, error) {
	return i.getDescriptor(model.CodeMediaType)
}

func (i *pulledImage) FetchFilter(ctx context.Context) (model.Filter, error) {
	desc, err := i.Descriptor()
	if err != nil {
		return nil, err
	}

	fetcher, err := i.resolver.Fetcher(ctx, i.ref)
	if err != nil {
		return nil, err
	}

	return fetcher.Fetch(ctx, desc)
}

func (i *pulledImage) FetchConfig(ctx context.Context) (*config.Config, error) {
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

func (i *pulledImage) getDescriptor(mediaType string) (ocispec.Descriptor, error) {
	for _, child := range i.children {
		if child.MediaType == mediaType {
			return child, nil
		}
	}
	return ocispec.Descriptor{}, errors.Errorf("media type %v not found on image", mediaType)
}
