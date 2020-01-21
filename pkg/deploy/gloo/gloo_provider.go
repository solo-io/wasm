package gloo

import (
	"context"

	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/wasm"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"
	corev1 "k8s.io/api/core/v1"
)

// selects the gateways to which to deploy the wasm filter(s)
type Selector struct {
	Namespaces    []string
	GatewayLabels map[string]string
}

type Provider struct {
	Ctx context.Context

	GatewayClient gatewayv1.GatewayClient

	// used to determine the workloads and gateways to which we apply the filters
	Selector Selector
}

// applies the filter to all selected workloads in selected namespaces
func (p *Provider) ApplyFilter(filter *v1.FilterSpec) error {
	return p.updateGateways(func(gateway *gatewayv1.Gateway) error {
		return apendWasmConfig(filter, gateway)
	})
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *v1.FilterSpec) error {
	return p.updateGateways(func(gateway *gatewayv1.Gateway) error {
		return removeWasmConfig(filter.Id, gateway)
	})
}
func (p *Provider) updateGateways(updateFunc func(gateway *gatewayv1.Gateway) error) error {
	namespaces := p.Selector.Namespaces
	if len(namespaces) == 0 {
		namespaces = []string{corev1.NamespaceAll}
	}
	for _, ns := range namespaces {
		gateways, err := p.GatewayClient.List(ns, clients.ListOpts{
			Ctx:      p.Ctx,
			Selector: p.Selector.GatewayLabels,
		})
		if err != nil {
			return err
		}

		for _, gw := range gateways {
			if err := updateFunc(gw); err != nil {
				contextutils.LoggerFrom(p.Ctx).Warnf("skipping gateway %v", gw.Metadata.Ref())
			}
			if _, err := p.GatewayClient.Write(gw, clients.WriteOpts{
				Ctx:               p.Ctx,
				OverwriteExisting: true,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: currently gloo only supports 1 wasm filter
// when it is updated, this should become an APPEND
func apendWasmConfig(filter *v1.FilterSpec, gateway *gatewayv1.Gateway) error {
	httpGw := gateway.GetHttpGateway()
	if httpGw == nil {
		return errors.Errorf("must contain httpGateway field")
	}
	opts := httpGw.GetOptions()
	if opts == nil {
		opts = &gloov1.HttpListenerOptions{}
		httpGw.Options = opts
	}

	// this SET should become an APPEND
	// when gloo supports multiple wasm filters
	opts.Wasm = &wasm.PluginSource{
		Image:  filter.Image,
		Config: filter.Config,
		Name:   filter.Id,
		RootId: filter.RootID,
		VmType: wasm.PluginSource_V8, // default to V8
	}

	return nil
}

// TODO: currently gloo only supports 1 wasm filter
// when it is updated, this should become a REMOVE
// currently it just sets to nil
func removeWasmConfig(filterID string, gateway *gatewayv1.Gateway) error {
	httpGw := gateway.GetHttpGateway()
	if httpGw == nil {
		return errors.Errorf("must contain httpGateway field")
	}
	opts := httpGw.GetOptions()

	if opts == nil {
		return nil
	}

	// when it is updated, this should become a REMOVE
	if opts.Wasm != nil && opts.Wasm.Name == filterID {
		opts.Wasm = nil
	}

	return nil
}
