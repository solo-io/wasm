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

	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"

	"github.com/solo-io/wasme/pkg/cache"
)

func watchFile(ctx context.Context, imageCache cache.Cache, refFile, directory string) error {
	// for each ref in the file, add it to the cache,
	// and if directory is not empty write it t here
	fw := fileWatcher{
		imageCache: imageCache,
		refFile:    refFile,
		directory:  directory,
	}

	return fw.watchFile(ctx)
}

type fileWatcher struct {
	imageCache cache.Cache
	refFile    string
	directory  string
}

func (f *fileWatcher) watchFileAndGetRefs(ctx context.Context, refFile string) <-chan string {
	res := make(chan string)
	go func() {
		defer close(res)
		for {
			var refs []string

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 10):
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

func (f *fileWatcher) watchFile(ctx context.Context) error {
	for ref := range f.watchFileAndGetRefs(ctx, f.refFile) {
		desc, err := f.imageCache.Add(ctx, ref)
		if err != nil {
			return err
		}
		err = f.addToDirectory(ctx, desc)
		if err != nil {
			return errors.Wrapf(err, "adding digest to directory %v", f.directory)
		}
	}
	return nil
}

func (f *fileWatcher) addToDirectory(ctx context.Context, digest digest.Digest) error {
	if f.directory == "" {
		return nil
	}
	// get filename from ref
	// check if filename exists
	filename := filepath.Join(f.directory, Digest2filename(digest))

	err := f.copyToFile(ctx, filename, digest)
	if err != nil {
	}
	return err
}

func (f *fileWatcher) copyToFile(ctx context.Context, filename string, digest digest.Digest) error {

	if _, err := os.Stat(filename); err == nil {
		// file already cached, nothing to do
		return nil
	}

	rc, err := f.imageCache.Get(ctx, digest)
	if err != nil {
		return err
	}
	defer rc.Close()

	// fail if file exists
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, rc)
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
