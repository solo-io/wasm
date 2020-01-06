---
title: "wasme deploy istio"
weight: 5
---
## wasme deploy istio

Deploy an Envoy WASM Filter to Istio Sidecar Proxies (Envoy).

### Synopsis

Deploy an Envoy WASM Filter to Istio Sidecar Proxies (Envoy).

wasme uses the EnvoyFilter Istio Custom Resource to pull and run wasm filters.
wasme deploys a server-side cache component which runs in cluster and pulls filter images.

Note: currently only Istio 1.4 is supported.


```
wasme deploy istio <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--labels <key1=val1,key2=val2>] [flags]
```

### Options

```
      --cache-custom-command strings   custom command to provide to the cache server image
      --cache-name string              name of resources for the wasm image cache server (default "wasme-cache")
      --cache-namespace string         namespace of resources for the wasm image cache server (default "wasme")
      --cache-repo string              name of the image repository to use for the cache server daemonset (default "quay.io/solo-io/wasme")
      --cache-tag string               image tag to use for the cache server daemonset (default "0.0.1")
  -h, --help                           help for istio
      --name string                    name of the deployment or daemonset into which to inject the filter. if not set, will apply to all workloads in the target namespace
  -n, --namespace string               namespace of the workload(s) to inject the filter. (default "default")
  -t, --workload-type string           type of workload into which the filter should be injected. possible values are deployment or daemonset (default "deployment")
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

