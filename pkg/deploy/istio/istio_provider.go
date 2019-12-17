package istio

import (
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/config"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	"github.com/solo-io/wasme/pkg/deploy"
	"github.com/solo-io/wasme/pkg/util"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// selects the istio proxies to which to deploy the wasm filter(s)
type Selector struct {
	Namespaces     []string
	WorkloadLabels map[string]string
	ListenerType   networkingv1alpha3.EnvoyFilter_PatchContext
}

type Provider struct {
	IstioClient versionedclient.Interface

	//global config namespace
	IstioConfigNamespace string

	// used to determine the workloads and listeners to which we apply the filters
	Selector Selector
}

// applies the filter to all selected workloads in selected namespaces
func (p *Provider) ApplyFilter(filter *deploy.Filter) error {
	namespaces := p.Selector.Namespaces
	if len(namespaces) == 0 {
		namespaces = []string{v1.NamespaceAll}
	}
	for _, ns := range namespaces {
		// store the write function so we can swap it with Create interchangeably
		write := p.IstioClient.NetworkingV1alpha3().EnvoyFilters(ns).Update

		// see if an envoyFilter CRD exists already for this filter
		envoyFilter, err := p.IstioClient.NetworkingV1alpha3().EnvoyFilters(ns).Get(filter.ID, metav1.GetOptions{})
		if err != nil {
			// ensure we write the filter to a valid namespace
			writeNamespace := ns
			if writeNamespace == v1.NamespaceAll {
				writeNamespace = p.IstioConfigNamespace
			}
			envoyFilter = &v1alpha3.EnvoyFilter{
				ObjectMeta: metav1.ObjectMeta{
					// in istio's case, filter ID must be a kube-compliant name
					Name:      filter.ID,
					Namespace: ns,
				},
			}

			// object does not exist so we must use crate
			write = p.IstioClient.NetworkingV1alpha3().EnvoyFilters(ns).Create
		}

		envoyFilter.Spec = makeSpec(filter, p.Selector.ListenerType, p.Selector.WorkloadLabels)

		// write the created/updated EnvoyFilter
		if _, err := write(envoyFilter); err != nil {
			return err
		}
	}
	return nil
}

// removes the filter from all selected workloads in selected namespaces
func (p *Provider) RemoveFilter(filter *deploy.Filter) error {
	namespaces := p.Selector.Namespaces
	if len(namespaces) == 0 {
		namespaces = []string{v1.NamespaceAll}
	}
	for _, ns := range namespaces {
		// delete the filter
		if err := p.IstioClient.NetworkingV1alpha3().EnvoyFilters(ns).Delete(filter.ID, nil); err != nil {
			return err
		}
	}
	return nil
}

// create the spec for the EnvoyFilter crd
func makeSpec(filter *deploy.Filter, listenerType networkingv1alpha3.EnvoyFilter_PatchContext, labels map[string]string) networkingv1alpha3.EnvoyFilter {
	filterCfg := &config.WasmService{
		Config: &config.PluginConfig{
			Name:          filter.ID,
			RootId:        filter.RootID,
			Configuration: filter.Config,
			VmConfig: &config.VmConfig{
				Runtime: "envoy.wasm.runtime.v8",
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

	marshalledConf, err := util.MarshalStruct(filterCfg)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	wasmFilterConfig := &envoyhttp.HttpFilter{
		Name: "envoy.filters.http.wasm",
		ConfigType: &envoyhttp.HttpFilter_Config{
			Config: marshalledConf,
		},
	}

	patchValue, err := protoutils.MarshalStruct(wasmFilterConfig)
	if err != nil {
		// this should NEVER HAPPEN!
		panic(err)
	}

	return networkingv1alpha3.EnvoyFilter{
		WorkloadSelector: &networkingv1alpha3.WorkloadSelector{
			Labels: labels,
		},
		ConfigPatches: []*networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{{
			ApplyTo: networkingv1alpha3.EnvoyFilter_HTTP_FILTER,
			Match: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
				Context: listenerType,
				ObjectTypes: &networkingv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
					Listener: &networkingv1alpha3.EnvoyFilter_ListenerMatch{
						FilterChain: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{
							Filter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
								Name: "envoy.http_connection_manager",
								SubFilter: &networkingv1alpha3.EnvoyFilter_ListenerMatch_SubFilterMatch{
									Name: "envoy.router",
								},
							},
						},
					},
				},
			},
			Patch: &networkingv1alpha3.EnvoyFilter_Patch{
				Operation: networkingv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
				Value: patchValue,
			},
		}},
	}
}

/*

 */
