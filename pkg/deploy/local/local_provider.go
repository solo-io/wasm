package local

import (
	"context"
	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoy_config_bootstrap_v2 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	envoy_config_filter_network_http_connection_manager_v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	"github.com/solo-io/wasme/pkg/deploy"
	envoyfilter "github.com/solo-io/wasme/pkg/deploy/filter"
	wasmeutil "github.com/solo-io/wasme/pkg/util"
	"io"
	"io/ioutil"
)

type Provider struct {
	Ctx context.Context

	// input config
	Input io.Reader

	// output config
	Output io.Writer

	// the destination for storing the filter on the local filesystem
	FilterPath string

	// Use JSON instead of YAML for config (defaults to false)
	UseJsonConfig bool
}

func (p *Provider) getConfig() (*envoy_config_bootstrap_v2.Bootstrap, error) {
	b, err := ioutil.ReadAll(p.Input)
	if err != nil {
		return nil, err
	}

	if !p.UseJsonConfig {
		b, err = yaml.YAMLToJSON(b)
		if err != nil {
			return nil, err
		}
	}

	var bootstrap envoy_config_bootstrap_v2.Bootstrap
	return &bootstrap, wasmeutil.UnmarshalBytes(b, &bootstrap)
}

func (p *Provider) writeConfig(bootstrap *envoy_config_bootstrap_v2.Bootstrap) error {
	b, err := wasmeutil.MarshalBytes(bootstrap)
	if err != nil {
		return err
	}

	if !p.UseJsonConfig {
		b, err = yaml.JSONToYAML(b)
		if err != nil {
			return err
		}
	}

	_, err = p.Output.Write(b)
	return err
}

// applies the filter to all selected workloads in selected namespaces
func (p *Provider) ApplyFilter(filter *deploy.Filter) error {
	cfg, err := p.getConfig()
	if err != nil {
		return err
	}

	if err := addFilterToListeners(filter, cfg.GetStaticResources().GetListeners(), p.FilterPath); err != nil {
		return err
	}

	return p.writeConfig(cfg)
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *deploy.Filter) error {
	cfg, err := p.getConfig()
	if err != nil {
		return err
	}

	if err := removeFilterFromListeners(filter, cfg.GetStaticResources().GetListeners()); err != nil {
		return err
	}

	return p.writeConfig(cfg)
}

// for each hcm in each filter (where it exists)
func forEachHcm(listeners []*envoy_api_v2.Listener, fn func(networkFilter *envoy_api_v2_listener.Filter, cfg *envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager) error) error {
	for _, listener := range listeners {
		for _, chain := range listener.GetFilterChains() {
			for _, networkFilter := range chain.GetFilters() {
				if networkFilter.GetName() == util.HTTPConnectionManager {
					var cfg envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager
					err := wasmeutil.UnmarshalStruct(networkFilter.GetConfig(), &cfg)
					if err != nil {
						return err
					}

					if err := fn(networkFilter, &cfg); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func addFilterToListeners(filter *deploy.Filter, listeners []*envoy_api_v2.Listener, filterPath string) error {

	wasmFilter := envoyfilter.MakeWasmFilter(filter, envoyfilter.MakeLocalDatasource(filterPath))

	return forEachHcm(listeners, func(networkFilter *envoy_api_v2_listener.Filter, cfg *envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager) error {
		for i, httpFilter := range cfg.GetHttpFilters() {
			if httpFilter.GetName() == wasmeutil.WasmFilterName {
				// if a wasm filter with the given id exists, return error
				var wasmFilterConfig config.WasmService
				err := wasmeutil.UnmarshalStruct(httpFilter.GetConfig(), cfg)
				if err != nil {
					return err
				}

				if wasmFilterConfig.GetConfig().GetName() == filter.ID {
					return errors.Errorf("filter with id %v already present", filter.ID)
				}
			}

			if httpFilter.GetName() == util.Router {
				// insert the filter before the router
				cfg.HttpFilters = append(cfg.HttpFilters[:i], wasmFilter, httpFilter)

				// update the HCM with our filter
				cfgStruct, err := wasmeutil.MarshalStruct(cfg)
				if err != nil {
					return err
				}

				networkFilter.ConfigType = &envoy_api_v2_listener.Filter_Config{
					Config: cfgStruct,
				}

				break
			}
		}
		return nil
	})
}

func removeFilterFromListeners(filter *deploy.Filter, listeners []*envoy_api_v2.Listener) error {
	return forEachHcm(listeners, func(networkFilter *envoy_api_v2_listener.Filter, cfg *envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager) error {
		for i, httpFilter := range cfg.GetHttpFilters() {
			if httpFilter.GetName() == wasmeutil.WasmFilterName {
				// if a wasm filter with the given id exists, return error
				var wasmFilterConfig config.WasmService
				err := wasmeutil.UnmarshalStruct(httpFilter.GetConfig(), cfg)
				if err != nil {
					return err
				}

				if wasmFilterConfig.GetConfig().GetName() == filter.ID {
					cfg.HttpFilters = append(cfg.HttpFilters[:i], cfg.HttpFilters[i+1:]...)
					break
				}
			}
		}
		return nil
	})
}
