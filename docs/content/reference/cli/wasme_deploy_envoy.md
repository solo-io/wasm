---
title: "wasme deploy envoy"
weight: 5
---
## wasme deploy envoy

Run Envoy locally in Docker and attach a WASM Filter.

### Synopsis


This command runs Envoy locally in docker using a static bootstrap configuration which includes 
the specified WASM filter image. 

The bootstrap can be generated from an internal default or a modified config provided by the user with --bootstrap.

The generated bootstrap config can be output to a file with --out. If using this option, Envoy will not be started locally.


```
wasme deploy envoy <image> [--config=<filter config>] [--bootstrap=<custom envoy bootstrap file>] [--envoy-image=<custom envoy image>] [flags]
```

### Options

```
  -b, --bootstrap wasme deploy envoy   Path to an Envoy bootstrap config. If set, wasme deploy envoy will run Envoy locally using the provided configuration file. Set -in=- to use stdin. If empty, will use a default configuration template with a single route to `jsonplaceholder.typicode.com`.
      --docker-run-args docker run     Set to provide additional args to the docker run command used to launch Envoy. Ignored if --out is set.
  -e, --envoy-image string             Name of the Docker image containing the Envoy binary (default "docker.io/istio/proxyv2:1.5.1")
      --envoy-run-args envoy           Set to provide additional args to the envoy command used to launch Envoy. Ignored if --out is set.
  -h, --help                           help for envoy
      --out string                     If set, write the modified Envoy configuration to this file instead of launching Envoy. Set -out=- to use stdout.
      --store string                   Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
```

### Options inherited from parent commands

```
      --config string    optional config that will be passed to the filter. accepts an inline string.
      --id string        unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
      --root-id string   optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
  -v, --verbose          verbose output
```

### SEE ALSO

* [wasme deploy](../wasme_deploy)	 - Deploy an Envoy WASM Filter to the data plane (Envoy proxies).

