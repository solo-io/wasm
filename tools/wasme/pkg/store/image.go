package store

import (
	"bytes"
	"context"
	"io/ioutil"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasm/tools/wasme/pkg/config"
	"github.com/solo-io/wasm/tools/wasme/pkg/model"
)

type Image interface {
	model.Image
}

// an image stored on disk
type storedImage struct {
	ref         string
	descriptor  ocispec.Descriptor
	filterBytes []byte
	config      *config.Runtime
}

func NewStorableImage(ref string, descriptor ocispec.Descriptor, filterBytes []byte, runtime *config.Runtime) (*storedImage, error) {
	ref, err := model.FullRef(ref)
	if err != nil {
		return nil, err
	}
	return &storedImage{
		ref:         ref,
		descriptor:  descriptor,
		filterBytes: filterBytes,
		config:      runtime,
	}, nil
}

func (i *storedImage) Ref() string {
	return i.ref
}

func (i *storedImage) Descriptor() (ocispec.Descriptor, error) {
	return i.descriptor, nil
}

func (i *storedImage) FetchFilter(ctx context.Context) (model.Filter, error) {
	filter := model.Filter(ioutil.NopCloser(bytes.NewBuffer(i.filterBytes)))
	return filter, nil
}

func (i *storedImage) FetchConfig(ctx context.Context) (*config.Runtime, error) {
	return i.config, nil
}
