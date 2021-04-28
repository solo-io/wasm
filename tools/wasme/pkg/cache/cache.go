package cache

import (
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/opencontainers/go-digest"
	"github.com/solo-io/wasm/tools/wasme/pkg/model"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
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

	logger *zap.SugaredLogger

	cacheState cacheState
}

func NewCache(puller pull.ImagePuller) Cache {
	return NewCacheWithConext(context.Background(), puller)
}

func NewCacheWithConext(ctx context.Context, puller pull.ImagePuller) Cache {
	return &CacheImpl{
		Puller: puller,
		logger: contextutils.LoggerFrom(ctx),
	}
}

func (c *CacheImpl) Add(ctx context.Context, ref string) (digest.Digest, error) {
	if img := c.cacheState.findImage(ref); img != nil {
		c.logger.Debugf("found cached image ref %v", ref)
		desc, err := img.Descriptor()
		if err != nil {
			return "", err
		}
		return desc.Digest, nil
	}

	c.logger.Debugf("attempting to pull image %v", ref)
	image, err := c.Puller.Pull(ctx, ref)
	if err != nil {
		return "", err
	}

	desc, err := image.Descriptor()
	if err != nil {
		return "", err
	}

	c.cacheState.add(image)

	c.logger.Debugf("pulled image %v (digest: %v)", ref, desc.Digest)

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
		ref := strings.TrimPrefix(r.URL.Path, "/")
		c.logger.Debugf("serving http request for image ref %v", ref)
		desc, err := c.Add(r.Context(), ref)
		if err != nil {
			c.logger.Errorf("failed to add or fetch descriptor %v: %v", ref, err)
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
		c.logger.Errorf("image with sha %v not found", sha)
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
		c.logger.Errorf("failed fetching image content")
		http.NotFound(rw, r)
		return
	}
	if closer, ok := filter.(io.ReadCloser); ok {
		defer closer.Close()
	}

	rw.Header().Set("Content-Type", desc.MediaType)
	rw.Header().Set("Etag", "\""+string(desc.Digest)+"\"")
	c.logger.Debugf("writing image content...")
	if rs, ok := filter.(io.ReadSeeker); ok {
		// content of digests never changes so set mod time to a constant
		// don't use zero time because serve content doesn't use that.
		modTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		http.ServeContent(rw, r, sha, modTime, rs)
	} else {
		c.logger.Debugf("writing image content")
		rw.Header().Add("Content-Length", strconv.Itoa(int(desc.Size)))
		if r.Method != "HEAD" {
			_, err = io.Copy(rw, filter)
			if err != nil {
				// TODO: use real log
				fmt.Printf("error http %v\n", err)
			}
		}
	}
	c.logger.Debugf("finished writing %v: %v bytes", image.Ref(), desc.Size)
}
