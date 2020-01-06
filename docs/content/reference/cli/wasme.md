---
title: "wasme"
weight: 5
---
## wasme



### Synopsis



### Options

```
  -c, --config stringArray   auth config path
  -h, --help                 help for wasme
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
```

### SEE ALSO

* [wasme build](../wasme_build)	 - Compile the filter to wasm using Bazel-in-Docker
* [wasme catalog](../wasme_catalog)	 - interact with catalog
* [wasme deploy](../wasme_deploy)	 - Deploy an Envoy WASM Filter to the data plane (Envoy proxies).
* [wasme init](../wasme_init)	 - Initialize a source directory for new Envoy WASM Filter.

The provided --base will determine the content of the created directory. The default is 
a C++ example filter compatible with the latest Envoy Wasm APIs.

Note that Istio 1.4 uses an older version of the Envoy Wasm APIs and users should 
use --base=cpp-istio to initialize a filter source directory for Istio.

* [wasme list](../wasme_list)	 - List Envoy WASM Filters published to webassemblyhub.io.
* [wasme login](../wasme_login)	 - login so you can push images to webassemblyhub.io and submit them to the curated catalog
* [wasme pull](../wasme_pull)	 - Pull wasm filters from remote registry
* [wasme push](../wasme_push)	 - Push wasm filter to remote registry
* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

