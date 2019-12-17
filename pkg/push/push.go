package push

import (
	"context"
	"fmt"
	"github.com/solo-io/gloo/pkg/utils/protoutils"
	"github.com/solo-io/wasme/pkg/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/containerd/containerd/reference"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/model"
)


type LocalFilter interface {
	CodeFilename() string
	DescriptorsFilename() string
	Image() string
	RootId() string
}

type localFilterImpl struct {
	codeFilename        string
	descriptorsFilename string
	image               string
	rootId              string
}

func NewLocalFilter(codeFilename, descriptorsFilename, image, rootId string) *localFilterImpl {
	return &localFilterImpl{
		codeFilename:        codeFilename,
		descriptorsFilename: descriptorsFilename,
		image:               image,
		rootId:              rootId,
	}
}

func (l *localFilterImpl) CodeFilename() string        { return l.codeFilename }
func (l *localFilterImpl) DescriptorsFilename() string { return l.descriptorsFilename }
func (l *localFilterImpl) Image() string               { return l.image }
func (l *localFilterImpl) RootId() string              { return l.rootId }

type Pusher interface {
	Push(ctx context.Context, f LocalFilter) error
}

type PusherImpl struct {
	Resolver   remotes.Resolver
	Authorizer docker.Authorizer
}

func makeConfig(localFilter LocalFilter) *config.Config {
	return &config.Config{
		Roots: []*config.Root{{Id: localFilter.RootId()}},
	}
}

func writeConfig(cfg *config.Config) (string, func(), error) {

	bytes, err := protoutils.MarshalBytes(cfg)
	if err != nil {
		return "", nil, err
	}

	tmpfile, err := ioutil.TempFile("", "config.*.txt")
	if err != nil {
		return "", nil, err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(bytes)
	if err != nil {
		// remove after close
		defer os.Remove(tmpfile.Name())
		return "", nil, err
	}

	return tmpfile.Name(), func() { os.Remove(tmpfile.Name()) }, nil
}

func (p *PusherImpl) Push(ctx context.Context, localFilter LocalFilter) error {
	var pushOpts []oras.PushOpt

	store := content.NewFileStore("")

	cfg := makeConfig(localFilter)
	filename, cleanup, err := writeConfig(cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	file, err := store.Add(config.ConfigFilename, model.ConfigMediaType, filename)
	if err != nil {
		return err
	}
	file.Annotations = nil
	pushOpts = append(pushOpts, oras.WithConfig(file))

	files, err := getFiles(localFilter, store)
	if err != nil {
		return err
	}

	p.checkAuth(ctx, localFilter.Image())

	desc, err := oras.Push(ctx, p.Resolver, localFilter.Image(), store, files, pushOpts...)
	if err != nil {
		return err
	}
	fmt.Println("Pushed", localFilter.Image())
	fmt.Println("Digest:", desc.Digest)

	return err
}

func (p *PusherImpl) checkAuth(ctx context.Context, ref string) {
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
		p.Authorizer.AddResponses(ctx, []*http.Response{resp})
	}
}

func getFiles(localFilter LocalFilter, store *content.FileStore) ([]ocispec.Descriptor, error) {
	var files []ocispec.Descriptor

	if descriptors := localFilter.DescriptorsFilename(); descriptors != "" { // TODO : error here instead?
		cfgFile, err := store.Add("config.proto.bin", model.ProtoSchemaMediaType, descriptors)
		if err != nil {
			return nil, err
		}
		files = append(files, cfgFile)
	}

	file, err := store.Add("code.wasm", model.CodeMediaType, localFilter.CodeFilename())
	if err != nil {
		return nil, err
	}
	files = append(files, file)
	return files, nil
}
