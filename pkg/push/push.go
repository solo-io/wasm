package push

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/extend-envoy/pkg/model"
)

const (
	annotationConfig   = "$config"
	annotationManifest = "$manifest"
)

type LocalFilter interface {
	CodeFilename() string
	ConfigFilename() string
	Image() string
}

type localFilterImpl struct {
	codeFilename   string
	configFilename string
	image          string
}

func NewLocalFilter(codeFilename, configFilename, image string) *localFilterImpl {
	return &localFilterImpl{
		codeFilename:   codeFilename,
		configFilename: configFilename,
		image:          image,
	}
}

func (l *localFilterImpl) CodeFilename() string   { return l.codeFilename }
func (l *localFilterImpl) ConfigFilename() string { return l.configFilename }
func (l *localFilterImpl) Image() string          { return l.image }

type Pusher interface {
	Push(ctx context.Context, f LocalFilter) error
}

type PusherImpl struct {
	Resolver remotes.Resolver
}

func (p *PusherImpl) Push(ctx context.Context, localFilter LocalFilter) error {
	var pushOpts []oras.PushOpt

	store := content.NewFileStore("")

	filename := localFilter.ConfigFilename()
	if filename != "" { // TODO : error here instead
		file, err := store.Add(annotationConfig, model.ConfigMediaType, filename)
		if err != nil {
			return err
		}
		file.Annotations = nil
		pushOpts = append(pushOpts, oras.WithConfig(file))
	}

	files, err := getFiles(localFilter, store)
	if err != nil {
		return err
	}

	desc, err := oras.Push(ctx, p.Resolver, localFilter.Image(), store, files, pushOpts...)
	if err != nil {
		return err
	}

	fmt.Println("Pushed", localFilter.Image())
	fmt.Println("Digest:", desc.Digest)

	return err
}

func getFiles(localFilter LocalFilter, store *content.FileStore) ([]ocispec.Descriptor, error) {

	var files []ocispec.Descriptor
	file, err := store.Add("code.wasm", model.CodeMediaType, localFilter.CodeFilename())
	if err != nil {
		return nil, err
	}
	files = append(files, file)
	return files, nil
}
