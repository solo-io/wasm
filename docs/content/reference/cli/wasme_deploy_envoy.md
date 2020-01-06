---
title: "wasme deploy envoy"
weight: 5
---
## wasme deploy envoy

Configure a local instance of Envoy to run a WASM Filter.

### Synopsis


Unlike `wasme deploy gloo` and `wasme deploy istio`, `wasme deploy envoy` only outputs the Envoy configuration required to run the filter with Envoy.

Launch Envoy using the output configuration to run the wasm filter.


```
wasme deploy envoy <image> --id=<unique id> [--config=<inline string>] [--root-id=<root id>] --in=<input config file> --out=<output config file> --filter <path to filter wasm> [--use-json] [flags]
```

### Options

```
  -f, --filter string   the path to the compiled filter wasm file. (default "filter.wasm")
  -h, --help            help for envoy
      --in string       the input configuration file. the filter config will be added to each listener found in the file. Set -in=- to use stdin. (default "envoy.yaml")
      --out string      the output configuration file. the resulting config will be written to the file. Set -out=- to use stdout. (default "envoy.yaml")
      --use-json        parse the input file as JSON instead of YAML
```

### Options inherited from parent commands

```
      --config string     optional config that will be passed to the filter. accepts an inline string.
      --dry-run           print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string         unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
      --insecure          allow connections to SSL registry without certs
  -p, --password string   registry password
      --plain-http        use plain http and not https
      --root-id string    optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
  -u, --username string   registry username
```

### SEE ALSO

* [wasme deploy](../wasme_deploy)	 - Deploy an Envoy WASM Filter to the data plane (Envoy proxies).

