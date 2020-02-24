---
title: "wasme undeploy gloo"
weight: 5
---
## wasme undeploy gloo

Remove an Envoy WASM Filter from the Gloo Gateway Proxies (Envoy).

### Synopsis

wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.


```
wasme undeploy gloo --id=<unique name> [flags]
```

### Options

```
      --config string           optional config that will be passed to the filter. accepts an inline string.
  -h, --help                    help for gloo
  -l, --labels stringToString   select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected. (default [])
  -n, --namespaces strings      deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.
      --root-id string          optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
```

### Options inherited from parent commands

```
      --dry-run     print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string   unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
  -v, --verbose     verbose output
```

### SEE ALSO

* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

