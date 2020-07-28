package local

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/solo-io/wasme/pkg/model"
	"github.com/solo-io/wasme/pkg/store"

	v1 "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoy_config_bootstrap_v2 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	envoy_config_filter_network_hcm_v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	envoyfilter "github.com/solo-io/wasme/pkg/deploy/filter"
	wasmeutil "github.com/solo-io/wasme/pkg/util"
)

type Runner struct {
	Ctx context.Context

	// input bootstrap config
	Input io.ReadCloser

	// output config YAML only (DryRun), rather than invoking docker run
	Output io.Writer

	// path to root storage dir
	// default is ~/.wasme/store
	Store store.Store

	// additional args passed to the `docker run` command when running Envoy. Ignored if using DryRyn
	DockerRunArgs []string

	// additional args passed to the `envoy` command when running Envoy in Docker. Ignored if using DryRyn
	EnvoyArgs []string

	// the image ref for Envoy to run with docker. Ignored if using DryRyn
	EnvoyDockerImage string
}

// applies the filter to all static listeners in the bootstrap config
func (p *Runner) RunFilter(filter *v1.FilterSpec) error {
	cfg, err := p.getConfig()
	if err != nil {
		return err
	}

	image, err := p.Store.Get(filter.Image)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve image. make sure to run `wasme pull %v` to pull the image to your local storage.", filter.Image)
	}
	if filter.RootID == "" {
		imageCfg, err := image.FetchConfig(p.Ctx)
		if err != nil {
			return err
		}
		roots := imageCfg.GetConfig().GetRootIds()
		if len(roots) == 0 {
			return errors.Errorf("found no root_id on image or in params")
		}

		// default to first root
		filter.RootID = roots[0]
	}

	// allow filter ID to be empty, as we don't care in local envoy
	if filter.Id == "" {
		filter.Id = filter.RootID
	}

	filterDir, err := p.Store.Dir(filter.Image)
	if err != nil {
		return err
	}
	filterDir, err = filepath.Abs(filterDir)
	if err != nil {
		return err
	}
	filterFile := filepath.Join(filterDir, model.CodeFilename)

	if err := addFilterToListeners(filter, cfg.GetStaticResources().GetListeners(), filterFile); err != nil {
		return err
	}

	configYaml, err := marshalConfig(cfg)
	if err != nil {
		return err
	}

	if p.Output != nil {
		_, err = p.Output.Write(configYaml)
		return err
	}

	logrus.Infof("mounting filter file at %v", filterFile)

	logrus.Debugf("using bootstrap config: \n%s", string(configYaml))

	ports, err := getListenerPorts(cfg)
	if err != nil {
		return err
	}

	dockerArgs := append([]string{
		"--rm",
		"--name", filter.Id,
		"--entrypoint", "envoy",
		"-v", filterDir + ":" + filterDir + ":ro",
	}, p.DockerRunArgs...)

	for _, port := range ports {
		dockerArgs = append(dockerArgs, "-p", fmt.Sprintf("%v:%v", port, port))
	}

	envoyArgs := append([]string{
		"--disable-hot-restart",
		"--config-yaml", string(configYaml),
	}, p.EnvoyArgs...)

	logrus.WithFields(logrus.Fields{
		"container_name": filter.Id,
		"envoy_image":    p.EnvoyDockerImage,
		"filter_image":   filter.Image,
	}).Infof("running envoy-in-docker")

	if err := wasmeutil.DockerRun(os.Stdout, os.Stderr, nil, p.EnvoyDockerImage, dockerArgs, envoyArgs); err != nil {
		return err
	}

	return nil
}

func (p *Runner) getConfig() (*envoy_config_bootstrap_v2.Bootstrap, error) {
	b, err := ioutil.ReadAll(p.Input)
	if err != nil {
		return nil, err
	}

	if err := p.Input.Close(); err != nil {
		return nil, err
	}

	b, err = yaml.YAMLToJSON(b)
	if err != nil {
		return nil, err
	}

	var bootstrap envoy_config_bootstrap_v2.Bootstrap
	return &bootstrap, wasmeutil.UnmarshalBytes(b, &bootstrap)
}

func marshalConfig(bootstrap *envoy_config_bootstrap_v2.Bootstrap) ([]byte, error) {
	b, err := wasmeutil.MarshalBytes(bootstrap)
	if err != nil {
		return nil, err
	}

	b, err = yaml.JSONToYAML(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getListenerPorts(bootstrap *envoy_config_bootstrap_v2.Bootstrap) ([]uint32, error) {
	var ports []uint32
	for _, listener := range bootstrap.GetStaticResources().GetListeners() {
		port := listener.GetAddress().GetSocketAddress().GetPortValue()
		if port != 0 {
			ports = append(ports, port)
		}
	}
	if port := bootstrap.GetAdmin().GetAddress().GetSocketAddress().GetPortValue(); port != 0 {
		ports = append(ports, port)
	}
	return ports, nil
}

// for each hcm in each filter (where it exists)
func forEachHcm(listeners []*envoy_api_v2.Listener, fn func(networkFilter *envoy_api_v2_listener.Filter, cfg *envoy_config_filter_network_hcm_v2.HttpConnectionManager) error) error {
	for _, listener := range listeners {
		for _, chain := range listener.GetFilterChains() {
			for _, networkFilter := range chain.GetFilters() {
				if networkFilter.GetName() == util.HTTPConnectionManager {
					var cfg envoy_config_filter_network_hcm_v2.HttpConnectionManager
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

func addFilterToListeners(filter *v1.FilterSpec, listeners []*envoy_api_v2.Listener, filterPath string) error {

	wasmFilter := envoyfilter.MakeIstioWasmFilter(filter, envoyfilter.MakeLocalDatasource(filterPath))

	return forEachHcm(listeners, func(networkFilter *envoy_api_v2_listener.Filter, cfg *envoy_config_filter_network_hcm_v2.HttpConnectionManager) error {
		for i, httpFilter := range cfg.GetHttpFilters() {
			if httpFilter.GetName() == wasmeutil.WasmFilterName {
				var wasmFilterConfig config.WasmService
				err := wasmeutil.UnmarshalStruct(httpFilter.GetConfig(), cfg)
				if err != nil {
					return err
				}

				if wasmFilterConfig.GetConfig().GetName() == filter.Id {
					return errors.Errorf("filter with id %v already present", filter.Id)
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

const DefaultEnvoyImage = "docker.io/istio/proxyv2:1.5.1"

const BasicEnvoyConfig = `
admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 8080 }
    filter_chains:
    - filters:
      - name: envoy.http_connection_manager
        config:
          codec_type: AUTO
          stat_prefix: ingress_http
          route_config:
            name: test
            virtual_hosts:
            - name: jsonplaceholder
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route:
                  cluster: static-cluster
                  auto_host_rewrite: true
          http_filters:
          - name: envoy.router
  clusters:
  - name: static-cluster
    connect_timeout: 0.25s
    type: LOGICAL_DNS
    lb_policy: ROUND_ROBIN
    dns_lookup_family: V4_ONLY
    tls_context:
      sni: jsonplaceholder.typicode.com
    hosts: [{ socket_address: { address: jsonplaceholder.typicode.com, port_value: 443, ipv4_compat: true } }]
`
