---
title: "wasme undeploy envoy"
weight: 5
---
## wasme undeploy envoy

Remove an Envoy WASM Filter from the Envoy listeners.

### Synopsis

wasme removes the deployed filter matching the given id. 


```
wasme undeploy envoy --id=<unique name> [flags]
```

### Options

```
  -h, --help         help for envoy
      --in string    the input configuration file. the filter config will be added to each listener found in the file. Set -in=- to use stdin. If empty, will use a default configuration template. (default "envoy.yaml")
      --out string   the output configuration file. the resulting config will be written to the file. Set -out=- to use stdout. (default "envoy.yaml")
      --use-json     parse the input file as JSON instead of YAML
```

### Options inherited from parent commands

```
  -d, --debug       debug mode
      --dry-run     print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string   unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
  -v, --verbose     verbose output
```

### SEE ALSO

* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

