package push

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/util"

	"github.com/solo-io/wasme/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/containerd/containerd/reference"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/model"
)

type Image = model.Image

type Pusher interface {
	Push(ctx context.Context, image Image) error
}

type pusher struct {
	resolver   remotes.Resolver
	authorizer docker.Authorizer
}

func NewPusher(resolver remotes.Resolver, authorizer docker.Authorizer) *pusher {
	return &pusher{resolver: resolver, authorizer: authorizer}
}

func (p *pusher) Push(ctx context.Context, image Image) error {
	return util.RetryOn500(func() error {
		return p.push(ctx, image)
	})
}

func (p *pusher) push(ctx context.Context, image Image) error {
	p.checkAuth(ctx, image.Ref())

	store := content.NewMemoryStore()

	cfg, err := image.FetchConfig(ctx)
	if err != nil {
		return err
	}

	cfgBytes, err := cfg.ToBytes()
	if err != nil {
		return err
	}

	cfgDescriptor := store.Add(model.ConfigFilename, model.ConfigMediaType, cfgBytes)

	filter, err := image.FetchFilter(ctx)
	if err != nil {
		return err
	}

	filterBytes, err := ioutil.ReadAll(filter)
	if err != nil {
		return err
	}

	filterDescriptor := store.Add(model.CodeFilename, model.ContentMediaType, filterBytes)

	files := []ocispec.Descriptor{
		cfgDescriptor,
		filterDescriptor,
	}

	annotations := ManifestAnnotations(cfg)

	imageDesciptor, err := oras.Push(ctx, p.resolver, image.Ref(), store, files,
		oras.WithConfig(cfgDescriptor),
		oras.WithManifestAnnotations(annotations),
	)
	if err != nil {
		return errors.Wrap(err, "oras push failed")
	}

	logrus.Infof("Pushed %v", image.Ref())
	logrus.Infof("Digest: %v", imageDesciptor.Digest)

	return err
}

func (p *pusher) checkAuth(ctx context.Context, ref string) {
	if p.authorizer == nil {
		return
	}
	refspec, err := reference.Parse(ref)
	if err != nil {
		return
	}
	url := url.URL{
		Host:   refspec.Hostname(),
		Path:   "/v2/",
		Scheme: "https",
	}
	if strings.HasPrefix(url.Host, "localhost:") || url.Host == "localhost" {
		url.Scheme = "http"
	}
	resp, err := http.Get(url.String())
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		p.authorizer.AddResponses(ctx, []*http.Response{resp})
	}
}

const (
	ManifestAnnotation_AbiVersion = "module.wasm.runtime/abi_version"
	ManifestAnnotation_Type       = "module.wasm.runtime/type"
)

func ManifestAnnotations(cfg *config.Runtime) map[string]string {
	return map[string]string{
		ManifestAnnotation_AbiVersion: strings.Join(cfg.AbiVersions, ","),
		ManifestAnnotation_Type:       cfg.Type,
	}
}
