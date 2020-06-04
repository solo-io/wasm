package cache

import (
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/opencontainers/go-digest"
	"github.com/solo-io/wasme/pkg/model"
	"github.com/solo-io/wasme/pkg/pull"
)

// Cache stores digests and image contents in memory
type Cache interface {
	// adds the image retrieve the digest for the image.
	// the digest will be cached if it is the initial Get
	Add(ctx context.Context, image string) (digest.Digest, error)

	// retrieve the wasm file from the image
	Get(ctx context.Context, digest digest.Digest) (model.Filter, error)
	http.Handler
}

type CacheImpl struct {
	Puller pull.ImagePuller

	cacheState cacheState
}

func NewCache(puller pull.ImagePuller) Cache {
	return &CacheImpl{
		Puller: puller,
	}
}

type cacheState struct {
	images     map[string]pull.Image
	imagesLock sync.RWMutex
}

func (c *cacheState) add(image pull.Image) {
	desc, err := image.Descriptor()
	if err != nil {
		// image is missing descriptor, should never happen
		// TODO: better logging impl
		log.Printf("error: image %v missing code descriptor", image.Ref())
		return
	}
	if c.find(desc.Digest) != nil {
		// check existence for idempotence
		// technically metadata can be different, but it's fine for now.
		return
	}
	c.imagesLock.Lock()
	if c.images == nil {
		c.images = make(map[string]pull.Image)
	}
	c.images[image.Ref()] = image
	c.imagesLock.Unlock()
}

func (c *cacheState) find(digest digest.Digest) pull.Image {
	c.imagesLock.RLock()
	defer c.imagesLock.RUnlock()
	if c.images == nil {
		return nil
	}
	for _, image := range c.images {
		desc, err := image.Descriptor()
		if err != nil {
			log.Printf("error: image %v missing code descriptor", image.Ref())
			return nil
		}

		if desc.Digest == digest {
			return image
		}
	}
	return nil
}
func (c *cacheState) findImage(image string) pull.Image {
	c.imagesLock.RLock()
	defer c.imagesLock.RUnlock()
	return c.images[image]
}

func (c *CacheImpl) Add(ctx context.Context, ref string) (digest.Digest, error) {
	if img := c.cacheState.findImage(ref); img != nil {
		desc, err := img.Descriptor()
		if err != nil {
			return "", err
		}
		return desc.Digest, nil
	}

	image, err := c.Puller.Pull(ctx, ref)
	if err != nil {
		return "", err
	}

	desc, err := image.Descriptor()
	if err != nil {
		return "", err
	}

	c.cacheState.add(image)

	return desc.Digest, nil
}

func (c *CacheImpl) Get(ctx context.Context, digest digest.Digest) (model.Filter, error) {
	image := c.cacheState.find(digest)
	if image == nil {
		return nil, errors.Errorf("image with digest %v not found", digest)
	}
	return image.FetchFilter(ctx)
}

func (c *CacheImpl) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// we support two paths:
	// /<HASH> - used in gloo
	// /image-name - used here to cache on demand
	_, file := path.Split(r.URL.Path)
	switch {
	case len(file) == hex.EncodedLen(crypto.SHA256.Size()):
		c.ServeHTTPSha(rw, r, file)
	default:
		// assume that the path is a ref. add it to cache
		desc, err := c.Add(r.Context(), strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		c.ServeHTTPSha(rw, r, desc.Encoded())
	}
}

func (c *CacheImpl) ServeHTTPSha(rw http.ResponseWriter, r *http.Request, sha string) {
	// parse the url
	ctx := r.Context()
	image := c.cacheState.find(digest.Digest("sha256:" + sha))
	if image == nil {
		http.NotFound(rw, r)
		return
	}

	desc, err := image.Descriptor()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	filter, err := image.FetchFilter(ctx)
	if err != nil {
		http.NotFound(rw, r)
		return
	}
	if closer, ok := filter.(io.ReadCloser); ok {
		defer closer.Close()
	}

	rw.Header().Set("Content-Type", desc.MediaType)
	rw.Header().Set("Etag", "\""+string(desc.Digest)+"\"")
	if rs, ok := filter.(io.ReadSeeker); ok {
		// content of digests never changes so set mod time to a constant
		// don't use zero time because serve content doesn't use that.
		modTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		http.ServeContent(rw, r, sha, modTime, rs)
	} else {
		rw.Header().Add("Content-Length", strconv.Itoa(int(desc.Size)))
		if r.Method != "HEAD" {
			_, err = io.Copy(rw, filter)
			if err != nil {
				// TODO: use real log
				fmt.Printf("error http %v\n", err)
			}
		}
	}
}
