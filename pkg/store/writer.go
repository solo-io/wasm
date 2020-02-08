package store

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/solo-io/wasme/pkg/model"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/solo-io/wasme/pkg/config"
)

const (
	// store the image ref for use in List() calls
	// TODO: consider making this generic metadata
	imageRefFilename   = "image_ref"
	descriptorFilename = "descriptor.json"
	configFilename     = model.ConfigFilename
	filterFilename     = model.CodeFilename
)

// writes an image into and reads an image out of a directory
type imageReadWriter struct {
	dir string
}

func (w imageReadWriter) writeRef(image Image) error {
	ref := image.Ref()
	imageRefFile := filepath.Join(w.dir, imageRefFilename)
	return ioutil.WriteFile(imageRefFile, []byte(ref), 0644)
}

func (w imageReadWriter) writeConfig(ctx context.Context, image Image) error {
	cfg, err := image.FetchConfig(ctx)
	if err != nil {
		return err
	}
	cfgBytes, err := cfg.ToBytes()
	if err != nil {
		return err
	}
	configFile := filepath.Join(w.dir, configFilename)
	return ioutil.WriteFile(configFile, cfgBytes, 0644)
}

func (w imageReadWriter) writeDescriptor(image Image) error {
	desc, err := image.Descriptor()
	if err != nil {
		return err
	}
	descBytes, err := json.Marshal(desc)
	if err != nil {
		return err
	}
	descriptorFile := filepath.Join(w.dir, descriptorFilename)
	return ioutil.WriteFile(descriptorFile, descBytes, 0644)
}

func (w imageReadWriter) writeFilter(ctx context.Context, image Image) error {
	filter, err := image.FetchFilter(ctx)
	if err != nil {
		return err
	}

	filterFile := filepath.Join(w.dir, filterFilename)
	destFile, err := os.Create(filterFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, filter)
	if err != nil {
		// to avoid partial copies, delete the dest file if it exists
		_ = os.Remove(filterFile)
		return err
	}
	return destFile.Close()
}

func (w imageReadWriter) writeImage(ctx context.Context, image Image) error {
	if err := os.MkdirAll(w.dir, 0777); err != nil {
		return err
	}

	if err := w.writeRef(image); err != nil {
		return err
	}
	if err := w.writeDescriptor(image); err != nil {
		return err
	}
	if err := w.writeConfig(ctx, image); err != nil {
		return err
	}
	if err := w.writeFilter(ctx, image); err != nil {
		return err
	}

	return nil
}

func (w imageReadWriter) readRef() (string, error) {
	imageRefFile := filepath.Join(w.dir, imageRefFilename)
	raw, err := ioutil.ReadFile(imageRefFile)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (w imageReadWriter) readDescriptor() (ocispec.Descriptor, error) {
	var desc ocispec.Descriptor

	descriptorFile := filepath.Join(w.dir, descriptorFilename)
	descBytes, err := ioutil.ReadFile(descriptorFile)
	if err != nil {
		return desc, err
	}

	return desc, json.Unmarshal(descBytes, &desc)
}

func (w imageReadWriter) readConfig() (*config.Runtime, error) {
	configFile := filepath.Join(w.dir, configFilename)
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg, err := config.FromBytes(raw)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (w imageReadWriter) readFilter() ([]byte, error) {
	filterFile := filepath.Join(w.dir, filterFilename)
	return ioutil.ReadFile(filterFile)
}

// will skip loading the filter
func (w imageReadWriter) readImage() (*storedImage, error) {
	ref, err := w.readRef()
	if err != nil {
		return nil, err
	}
	desc, err := w.readDescriptor()
	if err != nil {
		return nil, err
	}
	cfg, err := w.readConfig()
	if err != nil {
		return nil, err
	}
	filterBytes, err := w.readFilter()
	if err != nil {
		return nil, err
	}

	return NewStorableImage(ref, desc, filterBytes, cfg)
}
