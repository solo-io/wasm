---
title: "wasme deploy gloo"
weight: 5
---
## wasme deploy gloo

Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy).

### Synopsis

Deploys an Envoy WASM Filter to Gloo Gateway Proxies.

wasme uses the Gloo Gateway CR to pull and run wasm filters.

Use --namespaces to constrain the namespaces of Gateway CRs to update.

Use --labels to use a match Gateway CRs by label.


```
wasme deploy gloo <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>] [flags]
```

### Options

```
  -h, --help                    help for gloo
  -l, --labels stringToString   select deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected. (default [])
  -n, --namespaces strings      deploy the filter to selected Gateway resource in the given namespaces. if none provided, Gateways in all namespaces will be selected.
```

### Options inherited from parent commands

```
      --config string     optional config that will be passed to the filter. accepts an inline string.
  -d, --debug             debug mode
      --dry-run           print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string         unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
      --insecure          allow connections to SSL registry without certs
  -p, --password string   registry password
      --plain-http        use plain http and not https
      --root-id string    optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
  -u, --username string   registry username
  -v, --verbose           verbose output
```

### SEE ALSO

* [wasme deploy](../wasme_deploy)	 - Deploy an Envoy WASM Filter to the data plane (Envoy proxies).

