package filter

import (
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/util"
)

func MakeWasmFilter(filter *deploy.Filter) *envoyhttp.HttpFilter{
	filterCfg := &config.WasmService{
		Config: &config.PluginConfig{
			Name:          filter.ID,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			VmConfig: &config.VmConfig{
				Runtime: "envoy.wasm.runtime.v8", // default to v8
				Code: &core.AsyncDataSource{
					Specifier: &core.AsyncDataSource_Remote{
						Remote: &core.RemoteDataSource{
							HttpUri: &core.HttpUri{
								Uri: "TODO: URI",
								HttpUpstreamType: &core.HttpUri_Cluster{
									Cluster: "TODO: CLUSTER",
								},
								Timeout: &types.Duration{
									Seconds: 5, // TODO: customize
								},
							},
							Sha256: "TODO: SHA256",
						},
					},
				},
			},
		},
	}

	// here we need to use the golang proto marshal
	marshalledConf, err := util.MarshalStruct(filterCfg)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	return &envoyhttp.HttpFilter{
		Name: util.WasmFilterName,
		ConfigType: &envoyhttp.HttpFilter_Config{
			Config: marshalledConf,
		},
	}
}
