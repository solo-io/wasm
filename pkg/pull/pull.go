package pull

import (
	"context"

	"github.com/solo-io/wasme/pkg/model"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination ./mocks/pull.go github.com/solo-io/wasme/pkg/pull Image,ImageContent,ImagePuller

type Image = model.Image

// CodePuller Pulls the .wasm file descriptor from the image ref
// ImagePuller pulls oci image descriptors by their ref.
// Given the image as a wasme image,
// can also download files from those images
type ImagePuller interface {
	// Pull retrieves the child Descriptors for the provided image ref
	Pull(ctx context.Context, ref string) (Image, error)
}

type puller struct {
	resolver remotes.Resolver
}

func NewPuller(resolver remotes.Resolver) *puller {
	return &puller{
		resolver: resolver,
	}
}

func (p *puller) Pull(ctx context.Context, ref string) (Image, error) {
	ref = model.FullRef(ref)

	store := content.NewMemoryStore()

	name, desc, err := p.resolver.Resolve(ctx, ref)
	if err != nil {
		return nil, err
	}

	fetcher, err := p.resolver.Fetcher(ctx, ref)
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
	logrus.Debugf("%+v %+v %+v %+v\n", name, children, desc, err)

	return &pulledImage{
		children: children,
		ref:      ref,
		resolver: p.resolver,
	}, nil
}
