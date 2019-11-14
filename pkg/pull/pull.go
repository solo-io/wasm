package pull

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/extend-envoy/pkg/model"
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

type Puller interface {
	Pull(ctx context.Context, image string) (Filter, error)
	PullCodeDescriptor(ctx context.Context, ref string) (ocispec.Descriptor, error)
}
type PullerImpl struct {
	Resolver remotes.Resolver
}

func NewPuller(resolver remotes.Resolver) *PullerImpl {
	return &PullerImpl{
		Resolver: resolver,
	}
}

func (p *PullerImpl) PullCodeDescriptor(ctx context.Context, ref string) (ocispec.Descriptor, error) {

	store := content.NewMemoryStore()

	name, desc, err := p.Resolver.Resolve(ctx, ref)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	fetcher, err := p.Resolver.Fetcher(ctx, ref)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	_, err = remotes.FetchHandler(store, fetcher)(ctx, desc)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	children, err := images.ChildrenHandler(store)(ctx, desc)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	fmt.Printf("%+v %+v %+v %+v\n", name, children, desc, err)

	for _, child := range children {
		if child.MediaType == model.CodeMediaType {
			return child, nil
		}
	}
	return ocispec.Descriptor{}, errors.New("code not found")
}

func (p *PullerImpl) Pull(ctx context.Context, ref string) (Filter, error) {

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
