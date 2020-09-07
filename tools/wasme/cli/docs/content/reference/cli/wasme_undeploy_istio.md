---
title: "wasme undeploy istio"
weight: 5
---
## wasme undeploy istio

Remove an Envoy WASM Filter from the Istio Sidecar Proxies (Envoy).

### Synopsis

wasme uses the Istio EnvoyFilter CR to pull and run wasm filters.

Use --namespace to target workload(s) in a the namespaces of Gateway CRs to update.
Use --name to target a specific workload (deployment or daemonset) in the target namespace. If unspecified, all deployments 
in the namespace will be targeted.


```
wasme undeploy istio --id=<unique name> --namespace=<deployment namespace> [--name=<deployment name>] [flags]
```

### Options

```
      --cache-timeout duration   the length of time to wait for the server-side filter cache to pull the filter image before giving up with an error. set to 0 to skip the check entirely (note, this may produce a known race condition). (default 1m0s)
      --config string            optional config that will be passed to the filter. accepts an inline string.
  -h, --help                     help for istio
      --istio-namespace string   the namespace where the Istio control plane is installed (default "istio-system")
  -l, --labels stringToString    labels of the deployment or daemonset into which to inject the filter. if not set, will apply to all workloads in the target namespace (default [])
  -n, --namespace string         namespace of the workload(s) to inject the filter. (default "default")
      --root-id string           optional root ID used to bind the filter at the Envoy level. this value is normally read from the filter image directly.
  -t, --workload-type string     type of workload into which the filter should be injected. possible values are deployment or daemonset (default "deployment")
```

### Options inherited from parent commands

```
      --dry-run     print output any configuration changes to stdout rather than applying them to the target file / kubernetes cluster
      --id string   unique id for naming the deployed filter. this is used for logging as well as removing the filter. when running wasme deploy istio, this name must be a valid Kubernetes resource name.
  -v, --verbose     verbose output
```

### SEE ALSO

* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

