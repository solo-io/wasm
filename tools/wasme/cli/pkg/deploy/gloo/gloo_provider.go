package gloo

import (
	"context"
	"sort"

	skerrors "github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/wasm"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	v1 "github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1"
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
	return p.retryUpdateGateways(func(gateway *gatewayv1.Gateway) error {
		return apendWasmConfig(filter, gateway)
	})
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *v1.FilterSpec) error {
	return p.retryUpdateGateways(func(gateway *gatewayv1.Gateway) error {
		return removeWasmConfig(filter.Id, gateway)
	})
}

func (p *Provider) retryUpdateGateways(updateFunc func(gateway *gatewayv1.Gateway) error) error {
	return util.RetryOnFunc(func() error {
		return p.updateGateways(updateFunc)
	}, skerrors.IsResourceVersion)
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
			logrus.WithFields(logrus.Fields{
				"gateway": gw.Metadata.Namespace + "." + gw.Metadata.Name,
			}).Infof("updated gateway")
		}
	}
	return nil
}

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
	if opts.Wasm == nil {
		opts.Wasm = &wasm.PluginSource{}
	}

	filters := opts.Wasm.Filters

	// Gloo consumes the filter from this config
	glooWasmFilter := &wasm.WasmFilter{
		Image:  filter.Image,
		Config: filter.Config,
		Name:   filter.Id,
		RootId: filter.RootID,
		VmType: wasm.WasmFilter_V8, // default to V8
	}

	// see if a filter with this ID already exists, if so, update it
	var isUpdate bool
	for i, wasmFilter := range filters {
		if wasmFilter.Name == glooWasmFilter.Name {
			filters[i] = glooWasmFilter
			isUpdate = true
			logrus.WithFields(logrus.Fields{
				"filterID": glooWasmFilter.Name,
			}).Infof("updating wasm filter")
			break
		}
	}

	if !isUpdate {
		// append the user's filter
		logrus.WithFields(logrus.Fields{
			"filterID": glooWasmFilter.Name,
		}).Infof("appending wasm filter")
		filters = append(filters, glooWasmFilter)
	}

	// sort for idempotence
	sort.SliceStable(filters, func(i, j int) bool {
		return filters[i].Name < filters[j].Name
	})

	opts.Wasm.Filters = filters

	return nil
}

func removeWasmConfig(filterID string, gateway *gatewayv1.Gateway) error {
	httpGw := gateway.GetHttpGateway()
	if httpGw == nil {
		return errors.Errorf("must contain httpGateway field")
	}
	opts := httpGw.GetOptions()

	if opts == nil {
		return nil
	}

	// remove the filter
	filters := opts.GetWasm().GetFilters()
	for i, filter := range filters {
		if filter.GetName() == filterID {
			opts.Wasm.Filters = append(opts.Wasm.Filters[:i], opts.Wasm.Filters[i+1:]...)
			logrus.WithFields(logrus.Fields{
				"filterID": filterID,
			}).Infof("removing wasm filter")
			break
		}
	}

	return nil
}
