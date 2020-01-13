package filter

import (
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	wasmev1 "github.com/solo-io/wasme/operator/pkg/api/wasme.io/v1"
	"github.com/solo-io/wasme/pkg/util"
)

func MakeRemoteDataSource(uri, cluster string) *core.AsyncDataSource {
	return &core.AsyncDataSource{
		Specifier: &core.AsyncDataSource_Remote{
			Remote: &core.RemoteDataSource{
				HttpUri: &core.HttpUri{
					Uri: uri,
					HttpUpstreamType: &core.HttpUri_Cluster{
						Cluster: cluster,
					},
					Timeout: &types.Duration{
						Seconds: 5, // TODO: customize
					},
				},
			},
		},
	}
}

func MakeLocalDatasource(path string) *core.AsyncDataSource {
	return &core.AsyncDataSource{
		Specifier: &core.AsyncDataSource_Local{
			Local: &core.DataSource{
				Specifier: &core.DataSource_Filename{
					Filename: path,
				},
			},
		},
	}
}

func MakeFilenameDatasource(path string) *core.DataSource {
	return &core.DataSource{
		Specifier: &core.DataSource_Filename{
			Filename: path,
		},
	}
}

func MakeWasmFilter(filter *wasmev1.FilterSpec, dataSrc *core.AsyncDataSource) *envoyhttp.HttpFilter {
	filterCfg := &config.WasmService{
		Config: &config.PluginConfig{
			Name:          filter.Id,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			VmConfig: &config.VmConfig{
				Runtime: "envoy.wasm.runtime.v8", // default to v8
				Code:    dataSrc,
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

// istio wasm does nto support AsyncDataSource, we must manually stitch the struct
// TODO: remove this when Istio updates their Envoy APIs
func MakeHackyIstioWasmFilter(filter *wasmev1.FilterSpec, dataSrc *core.DataSource) *envoyhttp.HttpFilter {
	filterCfg := &config.WasmService{
		Config: &config.PluginConfig{
			Name:          filter.Id,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			VmConfig: &config.VmConfig{
				Runtime: "envoy.wasm.runtime.v8", // default to v8
			},
		},
	}

	// here we need to use the golang proto marshal
	marshalledConf, err := util.MarshalStruct(filterCfg)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	marshalledDataSrc, err := util.MarshalStruct(dataSrc)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	marshalledConf.Fields["config"].GetStructValue().Fields["vmConfig"].GetStructValue().Fields["code"] = &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: marshalledDataSrc}}

	return &envoyhttp.HttpFilter{
		Name: util.WasmFilterName,
		ConfigType: &envoyhttp.HttpFilter_Config{
			Config: marshalledConf,
		},
	}
}
