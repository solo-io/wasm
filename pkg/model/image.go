package model

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/deislabs/oras/pkg/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/util"
)

// a WASM Module Runtime
type Runtime struct {
	Type       RuntimeType
	AbiVersion string
}

type RuntimeType string

const (
	Runtime_EnvoyProxy RuntimeType = "envoy_proxy"
)

// represents the descriptors for an image, as well as accessors to the image contents
type Image interface {
	// ref to the image
	Ref() string

	// get the descriptor for the module layer
	Descriptor() (ocispec.Descriptor, error)

	// get the filter .wasm file from the image
	FetchFilter(ctx context.Context) (Filter, error)

	// get the filter config from the image
	FetchConfig(ctx context.Context) (*config.Runtime, error)
}

// media types stored in a Wasm Module image
const (
	ConfigMediaType  = "application/vnd.module.wasm.config.v1+json"
	ContentMediaType = "application/vnd.module.wasm.content.layer.v1+wasm"
)

// default filenames stored in a Wasm Module Image
const (
	ConfigFilename = "runtime-config.json"
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

	return store.Add(CodeFilename, ContentMediaType, bytes), nil
}

// expand the ref to contain :latest suffix if no tag provided
func FullRef(ref string) string {
	name, tag := util.SplitImageRef(ref)
	return name + ":" + tag
}
