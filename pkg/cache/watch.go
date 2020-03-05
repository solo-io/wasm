package cache

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// pulls images for a local cache
type LocalImagePuller interface {
	// watches ref file (images.txt) pulls each image to disk
	WatchFile(ctx context.Context) error
}

type localImagePuller struct {
	imageCache Cache
	refFile    string
	directory  string
}

func NewLocalImagePuller(imageCache Cache, refFile string, directory string) *localImagePuller {
	return &localImagePuller{imageCache: imageCache, refFile: refFile, directory: directory}
}

func (f *localImagePuller) WatchFile(ctx context.Context) error {
	logrus.Infof("starting writing images to %v, reading from %v", f.directory, f.refFile)
	for ref := range f.watchFileAndGetRefs(ctx, f.refFile) {
		logrus.Infof("pulling ref %v", ref)
		digest, err := f.imageCache.Add(ctx, ref)
		if err != nil {
			logrus.Infof("pulling error: %v", err)
			return err
		}
		err = f.addToDirectory(ctx, digest)
		if err != nil {
			logrus.Infof("writing image err: %v", err)
			return errors.Wrapf(err, "adding digest to directory %v", f.directory)
		}
	}
	return nil
}

func (f *localImagePuller) watchFileAndGetRefs(ctx context.Context, refFile string) <-chan string {
	res := make(chan string)
	go func() {
		defer close(res)
		for {
			var refs []string

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 2):
				var err error
				// read refs from file
				refs, err = fileToRefs(refFile)
				if err != nil {
					// TODO: log? panic?
					fmt.Fprintln(os.Stderr, "failed parsing refs file: "+err.Error())
				}
			}

			for _, ref := range refs {
				select {
				case <-ctx.Done():
					return
				case res <- ref:
				}
			}
		}
	}()
	return res
}

func (f *localImagePuller) addToDirectory(ctx context.Context, digest digest.Digest) error {
	if f.directory == "" {
		return nil
	}
	// get filename from ref
	// check if filename exists
	filename := filepath.Join(f.directory, Digest2filename(digest))

	logrus.Infof("writing image to %v", filename)

	err := f.copyToFile(ctx, filename, digest)
	if err != nil {
	}
	return err
}

func (f *localImagePuller) copyToFile(ctx context.Context, filename string, digest digest.Digest) error {

	if _, err := os.Stat(filename); err == nil {
		// file already cached, nothing to do
		return nil
	}

	filter, err := f.imageCache.Get(ctx, digest)
	if err != nil {
		return err
	}
	if closer, ok := filter.(io.ReadCloser); ok {
		defer closer.Close()
	}

	// fail if file exists
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, filter)
	file.Close()
	if err != nil {
		// to avoid partial copies, delete the file if it exists
		os.Remove(filename)
	}
	return err
}

func Digest2filename(digest digest.Digest) string {
	return digest.Encoded()
}

func fileToRefs(refFile string) ([]string, error) {
	var refs []string
	file, err := os.Open(refFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ref := strings.TrimSpace(scanner.Text())
		refs = append(refs, ref)
	}
	return refs, scanner.Err()
}
