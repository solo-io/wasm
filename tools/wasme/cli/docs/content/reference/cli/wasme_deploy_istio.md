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

If --name is not provided, all deployments in the targeted namespace will attach the filter.

Note: currently only Istio 1.5.x - 1.8.x are supported.


```
wasme deploy istio <image> --id=<unique name> [--config=<inline string>] [--root-id=<root id>] [--namespaces <comma separated namespaces>] [--name deployment-name] [flags]
```

### Options

```
      --cache-custom-command strings     custom command to provide to the cache server image
      --cache-image-pull-policy string   image pull policy for the cache server daemonset. see https://kubernetes.io/docs/concepts/containers/images/ (default "IfNotPresent")
      --cache-name string                name of resources for the wasm image cache server (default "wasme-cache")
      --cache-namespace string           namespace of resources for the wasm image cache server (default "wasme")
      --cache-repo string                name of the image repository to use for the cache server daemonset (default "quay.io/solo-io/wasme")
      --cache-tag string                 image tag to use for the cache server daemonset (default "dev")
      --cache-timeout duration           the length of time to wait for the server-side filter cache to pull the filter image before giving up with an error. set to 0 to skip the check entirely (note, this may produce a known race condition). (default 1m0s)
  -h, --help                             help for istio
      --ignore-version-check             set to disable abi version compatability check.
      --istio-namespace string           the namespace where the Istio control plane is installed (default "istio-system")
  -l, --labels stringToString            labels of the deployment or daemonset into which to inject the filter. if not set, will apply to all workloads in the target namespace (default [])
  -n, --namespace string                 namespace of the workload(s) to inject the filter. (default "default")
  -t, --workload-type string             type of workload into which the filter should be injected. possible values are daemonset, deployment, statefulset (default "deployment")
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

