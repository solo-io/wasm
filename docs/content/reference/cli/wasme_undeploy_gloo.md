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
  -h, --help                    help for gloo
  -l, --labels stringToString   select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected. (default [])
  -n, --namespaces strings      deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.
```

### Options inherited from parent commands

```
  -c, --config stringArray   auth config path
  -d, --debug                debug mode
      --dry-run              print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string            unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
  -v, --verbose              verbose output
```

### SEE ALSO

* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

