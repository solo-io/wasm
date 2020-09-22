package filter

import (
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	"github.com/solo-io/wasm/tools/wasme/pkg/util"

	udpav1 "github.com/cncf/udpa/go/udpa/type/v1"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	corev3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	wasmfiltersv3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/wasm/v3"
	wasmv3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/wasm/v3"
	wasmev1 "github.com/solo-io/wasm/tools/wasme/cli/pkg/operator/api/wasme.io/v1"
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

func MakeV3LocalDatasource(path string) *corev3.AsyncDataSource {
	return &corev3.AsyncDataSource{
		Specifier: &corev3.AsyncDataSource_Local{
			Local: &corev3.DataSource{
				Specifier: &corev3.DataSource_Filename{
					Filename: path,
				},
			},
		},
	}
}

// MakeWasmFilter creates wasm filters to be used with Envoy.
// This will also work with Gloo (but not Istio).
func MakeWasmFilter(filter *wasmev1.FilterSpec, dataSrc *corev3.AsyncDataSource) *envoyhttp.HttpFilter {
	filterCfg := &wasmv3.WasmService{
		Config: &wasmv3.PluginConfig{
			Name:          filter.Id,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			Vm: &wasmv3.PluginConfig_VmConfig{
				VmConfig: &wasmv3.VmConfig{
					Runtime: "envoy.wasm.runtime.v8", // default to v8
					Code:    dataSrc,
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

// MakeTypedIstioWasmFilter returns a wasm filter for use with Istio.
// This method works for versions of Istio 1.7+
func MakeTypedIstioWasmFilter(filter *wasmev1.FilterSpec, dataSrc *corev3.AsyncDataSource) (*envoyhttp.HttpFilter, error) {
	filterCfg := &wasmfiltersv3.Wasm{
		Config: &wasmv3.PluginConfig{
			Name:          filter.Id,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			Vm: &wasmv3.PluginConfig_VmConfig{
				VmConfig: &wasmv3.VmConfig{
					Runtime: "envoy.wasm.runtime.v8", // default to v8
					Code:    dataSrc,
					VmId:    filter.Id,
				},
			},
		},
	}

	// Need a structpb because of udpa, convert it to a typedStruct,
	// marshall it to any
	marshalledConf, err := util.MarshalStruct(filterCfg)
	if err != nil {
		return nil, err
	}
	typedStructConf := &udpav1.TypedStruct{
		TypeUrl: "type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm",
		Value:   marshalledConf,
	}

	value, err := proto.Marshal(typedStructConf)
	if err != nil {
		return nil, err
	}
	anyTypedConf := &any.Any{TypeUrl: "type.googleapis.com/udpa.type.v1.TypedStruct", Value: value}

	envoyFilter := &envoyhttp.HttpFilter{
		Name: util.WasmFilterName,
		ConfigType: &envoyhttp.HttpFilter_TypedConfig{
			TypedConfig: anyTypedConf,
		},
	}
	return envoyFilter, nil
}

// MakeIstioWasmFilter returns a wasm filter for use with Istio. This method only
// works for versions of Istio up to and including 1.6. It will soon be deprecated
func MakeIstioWasmFilter(filter *wasmev1.FilterSpec, dataSrc *core.AsyncDataSource) (*envoyhttp.HttpFilter, error) {
	var cfgVal string
	if filter.Config != nil {
		// As the config's value is a StringValue, we need to unmarshall it,
		// typecheck it, then get the value out of the result.
		var da types.DynamicAny
		if err := types.UnmarshalAny(filter.Config, &da); err != nil {
			return nil, err
		}

		cfg, ok := da.Message.(*types.StringValue)
		if !ok {
			return nil, errors.Errorf("wasm filter configuration has an invalid type, should be StringValue")
		}
		cfgVal = cfg.GetValue()
	}

	filterCfg := &config.WasmService{
		Config: &config.PluginConfig{
			Name:          filter.Id,
			RootId:        filter.RootID,
			Configuration: cfgVal,
			VmConfig: &config.VmConfig{
				Runtime: "envoy.wasm.runtime.v8", // default to v8
				Code:    dataSrc,
				VmId:    filter.Id,
			},
		},
	}

	// here we need to use the golang proto marshal
	marshalledConf, err := util.MarshalStruct(filterCfg)
	if err != nil {
		return nil, err
	}

	return &envoyhttp.HttpFilter{
		Name: util.WasmFilterName,
		ConfigType: &envoyhttp.HttpFilter_Config{
			Config: marshalledConf,
		},
	}, nil
}
