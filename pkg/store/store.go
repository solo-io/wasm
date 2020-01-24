package store

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-multierror"
)

// a filter image that can be stored as a directory on-disk

type Store interface {
	List() ([]Image, error)
	Add(ctx context.Context, image Image) error
	Get(ref string) (*storedImage, error)
	Delete(ref string) error
}

type store struct {
	storageDir string
}

var defaultStorageDir = os.Getenv("HOME") + "/.wasme/store"

func NewStore(storageDir string) *store {
	if storageDir == "" {
		storageDir = defaultStorageDir
	}
	return &store{storageDir: storageDir}
}

// for the sake of efficiency, listing images
// does NOT load the wasm filter
// use Get() on the image ref to load the image
func (s *store) List() ([]Image, error) {
	files, err := ioutil.ReadDir(s.storageDir)
	if err != nil {
		return nil, err
	}

	var images []Image
	var readErrors error
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		image, err := imageReadWriter{dir: file.Name()}.readImage()
		if err != nil {
			readErrors = multierror.Append(readErrors, err)
			continue
		}

		images = append(images, image)
	}

	return images, readErrors
}

func (s *store) Add(ctx context.Context, image Image) error {
	return imageReadWriter{dir: dirname(image.Ref())}.writeImage(ctx, image)
}

func (s *store) Get(ref string) (*storedImage, error) {
	dir := dirname(ref)
	return imageReadWriter{dir: dir}.readImage()
}

func (s *store) Delete(ref string) error {
	dir := dirname(ref)
	return os.RemoveAll(dir)
}

func dirname(ref string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(ref)))
}
