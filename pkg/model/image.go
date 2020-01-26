package model

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"github.com/deislabs/oras/pkg/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/config"
)

// represents the descriptors for an image, as well as accessors to the image contents
type Image interface {
	// ref to the image
	Ref() string

	// get the image (Filter) descriptor
	Descriptor() (ocispec.Descriptor, error)

	// get the filter .wasm file from the image
	FetchFilter(ctx context.Context) (Filter, error)

	// get the filter config from the image
	FetchConfig(ctx context.Context) (*config.Config, error)
}

// media types stored in a Wasm Module image
const (
	ConfigMediaType = "application/vnd.io.solo.wasm.config.v1+json"
	CodeMediaType   = "application/vnd.io.solo.wasm.code.v1+wasm"
)

// default filenames stored in a Wasm Module Image
const (
	ConfigFilename = "config.json"
	CodeFilename   = "filter.wasm"
)

// a reader with access to the filter code
type Filter io.Reader

// helper function to get the descriptor for a wasm binary
func GetDescriptor(filter Filter) (ocispec.Descriptor, error) {
	store := content.NewMemoryStore()

	bytes, err := ioutil.ReadAll(filter)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return store.Add(CodeFilename, CodeMediaType, bytes), nil
}

// expand the ref to contain :latest suffix if no tag provided
func FullRef(ref string) string {
	parts := strings.Split(ref, ":")
	if len(parts) == 2 {
		return ref
	}
	return ref + ":latest"
}
