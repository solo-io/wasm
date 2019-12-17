package deploy

import (
	"context"
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/pull"
	v1 "k8s.io/api/core/v1"
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

// inspector inspects a filter
type Inspector interface {
	// gets the default root id for the filter. if
	// more than 1 root id is defined, defaults to the first
	SetRootID(f *Filter) error
}

type inspector struct {
	ctx    context.Context
	puller pull.Puller
}

func NewInspector(ctx context.Context, puller pull.Puller) *inspector {
	return &inspector{ctx: ctx, puller: puller}
}

// gets the root ID of the filter.
// the first time it must pull the image and inspect it
// second time it will cache it locally
// if the user provides
func (d *inspector) SetRootID(f *Filter) error {
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
func (d *inspector) getRootId(image string) (string, error) {
	cfg, err := d.puller.PullConfigFile(d.ctx, image)
	if err != nil {
		return "", err
	}
	if len(cfg.Roots) < 1 {
		return "", errors.Errorf("no roots found in config")
	}
	return cfg.Roots[0].Id, nil
}

// FilterDeployer deploys a wasm filter
// to one or more Envoy instances
// Each control plane provider has
// its own implementation
type FilterDeployer interface {
	// deploy the filter with initial config
	Deploy(filter Filter) error

	// update the filter config
	UpdateConfig(filter Filter) error

	// remove the filter
	Undeploy(filter Filter) error
}

// a Target for adding filters.
type Target v1.ObjectReference

// providers provide the interface to interact with a control plane
type Provider interface {
	GetFilters() ([]*Filter, error)
	SetFilters([]*Filter) error
}

func DeployFilter(filter *Filter, inspector Inspector, p Provider) error {
	// set rootID if not set by the user
	if err := inspector.SetRootID(filter); err != nil {
		return err
	}
	filters, err := p.GetFilters()
	if err != nil {
		return err
	}
	for _, f := range filters {
		if f.ID == filter.ID {
			return errors.Errorf("filter with id %v already exists. try with a different ID")
		}
	}
	filters = append(filters, filter)

	return p.SetFilters(filters)
}
