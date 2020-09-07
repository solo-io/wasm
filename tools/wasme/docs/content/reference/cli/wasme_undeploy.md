---
title: "wasme undeploy"
weight: 5
---
## wasme undeploy

Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

### Synopsis

Removes a deployed Envoy WASM Filter from Envoy instances.



### Options

```
      --dry-run     print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
  -h, --help        help for undeploy
      --id string   unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - The tool for building, pushing, and deploying Envoy WebAssembly Filters
* [wasme undeploy gloo](../wasme_undeploy_gloo)	 - Remove an Envoy WASM Filter from the Gloo Gateway Proxies (Envoy).
* [wasme undeploy istio](../wasme_undeploy_istio)	 - Remove an Envoy WASM Filter from the Istio Sidecar Proxies (Envoy).

