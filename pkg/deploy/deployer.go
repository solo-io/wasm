package deploy

import (
	"context"

	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"

	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/pull"
)

//go:generate mockgen -destination ./mocks/deployer.go github.com/solo-io/wasme/pkg/deploy Provider

// mesh-provider specific implementation that adds/removes filters
type Provider interface {
	ApplyFilter(filter *v1.FilterSpec) error
	RemoveFilter(filter *v1.FilterSpec) error
}

type Deployer struct {
	Ctx      context.Context
	Puller   pull.ImagePuller
	Provider Provider
}

func (d *Deployer) ApplyFilter(filter *v1.FilterSpec) error {
	if err := d.setRootID(filter); err != nil {
		return err
	}
	return d.Provider.ApplyFilter(filter)
}

func (d *Deployer) RemoveFilter(filter *v1.FilterSpec) error {
	return d.Provider.RemoveFilter(filter)
}

// gets the root ID of the filter.
// the first time it must pull the image and inspect it
// second time it will cache it locally
// if the user provides
func (d *Deployer) setRootID(f *v1.FilterSpec) error {
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
func (d *Deployer) getRootId(ref string) (string, error) {
	image, err := d.Puller.Pull(d.Ctx, ref)
	if err != nil {
		return "", err
	}

	cfg, err := image.FetchConfig(d.Ctx)
	if err != nil {
		return "", err
	}
	if len(cfg.Roots) < 1 {
		return "", errors.Errorf("no roots found in config")
	}
	return cfg.Roots[0].Id, nil
}
