package pull

import (
	"context"

	"github.com/solo-io/wasm/tools/wasme/pkg/util"

	"github.com/solo-io/wasm/tools/wasme/pkg/model"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination ./mocks/pull.go github.com/solo-io/wasm/tools/wasme/pkg/pull ImagePuller

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
	var image Image
	err := util.RetryOn500(func() error {
		var err error
		image, err = p.pull(ctx, ref)
		return err
	})
	return image, err
}

func (p *puller) pull(ctx context.Context, ref string) (Image, error) {
	ref, err := model.FullRef(ref)
	if err != nil {
		return nil, err
	}

	store := content.NewMemoryStore()

	name, manifest, err := p.resolver.Resolve(ctx, ref)
	if err != nil {
		return nil, err
	}

	fetcher, err := p.resolver.Fetcher(ctx, ref)
	if err != nil {
		return nil, err
	}
	_, err = remotes.FetchHandler(store, fetcher)(ctx, manifest)
	if err != nil {
		return nil, err
	}

	children, err := images.ChildrenHandler(store)(ctx, manifest)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("%+v %+v %+v %+v\n", name, children, manifest, err)

	return &pulledImage{
		children: children,
		ref:      ref,
		resolver: p.resolver,
	}, nil
}
