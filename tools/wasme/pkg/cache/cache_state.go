package cache

import (
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasm/tools/wasme/pkg/pull"
	"sync"
)

type cacheState struct {
	images     map[string]pull.Image
	imagesLock sync.RWMutex
}

func (c *cacheState) add(image pull.Image) {
	desc, err := image.Descriptor()
	if err != nil {
		// image is missing descriptor, should never happen
		// TODO: better logging impl
		logrus.Errorf("error: image %v missing code descriptor", image.Ref())
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
	logrus.Debugf("added image " + desc.Digest.String())
	c.imagesLock.Unlock()
}

func (c *cacheState) find(digest digest.Digest) pull.Image {
	c.imagesLock.RLock()
	defer c.imagesLock.RUnlock()
	if c.images == nil {
		return nil
	}
	logrus.Debugf("searching for image " + digest.String())
	for _, image := range c.images {
		desc, err := image.Descriptor()
		if err != nil {
			logrus.Errorf("error: image %v missing code descriptor", image.Ref())
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

func (c *cacheState) remove(image string) {
	c.imagesLock.Lock()
	defer c.imagesLock.Unlock()
	delete(c.images, image)
}
