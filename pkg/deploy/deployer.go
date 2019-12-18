package deploy

import (
	"context"
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/pull"
)

// the filter to deploy
type Filter struct {
	// unique identifier that will be used
	// to remove the filter as well as for logging
	ID string

	// name of image which houses the compiled wasm filter
	Image string

	// string of the config sent to the wasm filter
	// Currently has to be json or will crash
	Config string

	// the root id must match the root id
	// defined inside the filter.
	// if the user does not provide this field,
	// wasme will attempt to pull the image
	// and set it from the filter_conf
	// the first time it must pull the image and inspect it
	// second time it will cache it locally
	// if the user provides
	RootID string
}

// mesh-provider specific implementation that adds/removes filters
type Provider interface {
	ApplyFilter(filter *Filter) error
	RemoveFilter(filter *Filter) error
}

type Deployer struct {
	Ctx      context.Context
	Puller   pull.Puller
	Provider Provider
}

func (d *Deployer) ApplyFilter(filter *Filter) error {
	if err := d.setRootID(filter); err != nil {
		return err
	}
	return d.Provider.ApplyFilter(filter)
}

func (d *Deployer) RemoveFilter(filter *Filter) error {
	if err := d.setRootID(filter); err != nil {
		return err
	}
	return d.Provider.RemoveFilter(filter)
}

// gets the root ID of the filter.
// the first time it must pull the image and inspect it
// second time it will cache it locally
// if the user provides
func (d *Deployer) setRootID(f *Filter) error {
	if f.RootID != "" {
		return nil
	}
	rootId, err := d.getRootId(f.Image)
	if err != nil {
		return err
	}
	f.RootID = rootId
	return nil
}

// get the root id by pulling the image
func (d *Deployer) getRootId(image string) (string, error) {
	cfg, err := d.Puller.PullConfigFile(d.Ctx, image)
	if err != nil {
		return "", err
	}
	if len(cfg.Roots) < 1 {
		return "", errors.Errorf("no roots found in config")
	}
	return cfg.Roots[0].Id, nil
}
