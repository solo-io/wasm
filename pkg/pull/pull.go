package pull

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/config"
	"github.com/solo-io/wasme/pkg/model"
)

type FilterImpl struct {
	code io.ReadCloser
}

func (f *FilterImpl) Code() io.ReadCloser {
	return f.code
}

func (f *FilterImpl) Configs() []FilterConfig {
	return nil
}

// CodePuller Pulls the .wasm file descriptor from the image ref
type CodePuller interface {
	PullCodeDescriptor(ctx context.Context, ref string) (ocispec.Descriptor, error)
}

// ImagePuller pulls oci image descriptors by their ref.
// Given the image as a wasme image,
// can also download files from those images
type ImagePuller interface {
	CodePuller
	Pull(ctx context.Context, ref string) ([]ocispec.Descriptor, error)
	PullFilter(ctx context.Context, image string) (Filter, error)
	PullConfigFile(ctx context.Context, ref string) (*config.Config, error)
}

type PullerImpl struct {
	Resolver remotes.Resolver
}

func NewPuller(resolver remotes.Resolver) *PullerImpl {
	return &PullerImpl{
		Resolver: resolver,
	}
}

func (p *PullerImpl) Pull(ctx context.Context, ref string) ([]ocispec.Descriptor, error) {

	store := content.NewMemoryStore()

	name, desc, err := p.Resolver.Resolve(ctx, ref)
	if err != nil {
		return nil, err
	}

	fetcher, err := p.Resolver.Fetcher(ctx, ref)
	if err != nil {
		return nil, err
	}
	_, err = remotes.FetchHandler(store, fetcher)(ctx, desc)
	if err != nil {
		return nil, err
	}

	children, err := images.ChildrenHandler(store)(ctx, desc)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v %+v %+v %+v\n", name, children, desc, err)
	return children, nil
}

func (p *PullerImpl) PullCodeDescriptor(ctx context.Context, ref string) (ocispec.Descriptor, error) {

	children, err := p.Pull(ctx, ref)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	for _, child := range children {
		if child.MediaType == model.CodeMediaType {
			return child, nil
		}
	}
	return ocispec.Descriptor{}, errors.New("code not found")
}

func (p *PullerImpl) PullConfigFile(ctx context.Context, ref string) (*config.Config, error) {
	children, err := p.Pull(ctx, ref)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		if child.MediaType == model.ConfigMediaType {

			fetcher, err := p.Resolver.Fetcher(ctx, ref)
			if err != nil {
				return nil, err
			}

			rc, err := fetcher.Fetch(ctx, child)
			if err != nil {
				return nil, err
			}

			b, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			return config.FromBytes(b)
		}
	}
	return nil, errors.New("config not found")
}

func (p *PullerImpl) PullFilter(ctx context.Context, ref string) (Filter, error) {

	desc, err := p.PullCodeDescriptor(ctx, ref)
	if err != nil {
		return nil, err
	}
	return p.Fetch(ctx, ref, desc)
}

func (p *PullerImpl) Fetch(ctx context.Context, ref string, desc ocispec.Descriptor) (Filter, error) {

	fetcher, err := p.Resolver.Fetcher(ctx, ref)
	if err != nil {
		return nil, err
	}

	rc, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}
	return &FilterImpl{code: rc}, nil
}
