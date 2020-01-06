---
title: "wasme deploy"
weight: 5
---
## wasme deploy

Deploy an Envoy WASM Filter to the data plane (Envoy proxies).

### Synopsis

Deploys an Envoy WASM Filter to Envoy instances.

You must provide a value for --id which will become the unique ID of the deployed filter. When using --provider=istio, the ID must be a valid Kubernetes resource name.

You must specify --root-id unless a default root id is provided in the image configuration. Use --root-id to select the filter to run if the wasm image contains more than one filter.



### Options

```
      --dry-run          print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
  -h, --help             help for deploy
      --id string        unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
      --root-id string   optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
```

### Options inherited from parent commands

```
  -c, --config stringArray   auth config path
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
```

### SEE ALSO

* [wasme](../wasme)	 - 
* [wasme deploy envoy](../wasme_deploy_envoy)	 - Configure a local instance of Envoy to run a WASM Filter.
* [wasme deploy gloo](../wasme_deploy_gloo)	 - Deploy an Envoy WASM Filter to the Gloo Gateway Proxies (Envoy).
* [wasme deploy istio](../wasme_deploy_istio)	 - Deploy an Envoy WASM Filter to Istio Sidecar Proxies (Envoy).

