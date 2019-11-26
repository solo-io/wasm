package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/containerd/containerd/remotes"
	"github.com/solo-io/extend-envoy/pkg/pull"
)

type Cache interface {
	Add(ctx context.Context, image string) (digest.Digest, error)
	Get(ctx context.Context, digest digest.Digest) (io.ReadCloser, error)
	http.Handler
}

type CacheImpl struct {
	Puller   pull.Puller
	Resolver remotes.Resolver

	cacheState cacheState
}

type fetchableDescriptor struct {
	ocispec.Descriptor
	fetcher func(ctx context.Context) (io.ReadCloser, error) //remotes.Fetcher
}

type cacheState struct {
	descriptors     []fetchableDescriptor
	descriptorsLock sync.RWMutex
}

func (c *cacheState) add(d fetchableDescriptor) {
	if c.find(d.Digest) != nil {
		// check existance for idempotency
		// technically metadata can be different, but its fine for now.
		return
	}
	c.descriptorsLock.Lock()
	c.descriptors = append(c.descriptors, d)
	c.descriptorsLock.Unlock()
}

func (c *cacheState) find(digest digest.Digest) *fetchableDescriptor {
	c.descriptorsLock.RLock()
	defer c.descriptorsLock.RUnlock()
	for _, d := range c.descriptors {
		if d.Digest == digest {
			d := d
			return &d
		}
	}
	return nil
}

func (c *CacheImpl) Add(ctx context.Context, image string) (digest.Digest, error) {
	desc, err := c.Puller.PullCodeDescriptor(ctx, image)
	if err != nil {
		return digest.Digest(""), err
	}

	fd := fetchableDescriptor{
		Descriptor: desc,
		fetcher: func(subctx context.Context) (io.ReadCloser, error) {
			fetcher, err := c.Resolver.Fetcher(subctx, image)
			if err != nil {
				return nil, err
			}
			return fetcher.Fetch(subctx, desc)
		},
	}

	c.cacheState.add(fd)

	return desc.Digest, err
}

func (c *CacheImpl) Get(ctx context.Context, digest digest.Digest) (io.ReadCloser, error) {
	desc := c.cacheState.find(digest)
	if desc == nil {
		return nil, errors.New("not found")
	}
	return desc.fetcher(ctx)
}

func (c *CacheImpl) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// parse the url
	ctx := r.Context()
	_, file := path.Split(r.URL.Path)
	desc := c.cacheState.find(digest.Digest("sha256:" + file))
	if desc == nil {
		http.NotFound(rw, r)
		return
	}

	rc, err := desc.fetcher(ctx)
	if err != nil {
		http.NotFound(rw, r)
		return
	}
	defer rc.Close()

	rw.Header().Set("Content-Type", desc.MediaType)
	rw.Header().Set("Etag", "\""+string(desc.Digest)+"\"")
	if rs, ok := rc.(io.ReadSeeker); ok {
		// content of digests never changes so set mod time to a constant
		// don't use zero time because serve content doesn't use that.
		modTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		http.ServeContent(rw, r, file, modTime, rs)
	} else {
		rw.Header().Add("Content-Length", strconv.Itoa(int(desc.Size)))
		if r.Method != "HEAD" {
			_, err = io.Copy(rw, rc)
			if err != nil {
				// TODO: use real log
				fmt.Printf("error http %v\n", err)
			}
		}
	}
}
