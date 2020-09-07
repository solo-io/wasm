package pull

import (
	"context"
	"io"

	"github.com/solo-io/wasm/tools/wasme/pkg/model"

	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/solo-io/wasm/tools/wasme/pkg/config"
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
	return i.getDescriptor(model.ContentMediaType)
}

func (i *pulledImage) FetchFilter(ctx context.Context) (model.Filter, error) {
	desc, err := i.Descriptor()
	if err != nil {
		return nil, err
	}

	return i.fetchBlob(ctx, desc)
}

func (i *pulledImage) FetchConfig(ctx context.Context) (*config.Runtime, error) {
	desc, err := i.getDescriptor(model.ConfigMediaType)
	if err != nil {
		return nil, err
	}

	rc, err := i.fetchBlob(ctx, desc)
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

func (i *pulledImage) fetchBlob(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {
	fetcher, err := i.resolver.Fetcher(ctx, i.ref)
	if err != nil {
		return nil, err
	}

	return fetcher.Fetch(ctx, desc)
}
