---
title: "wasme"
weight: 5
---
## wasme



### Synopsis



### Options

```
  -c, --config stringArray   auth config path
  -d, --debug                debug mode
  -h, --help                 help for wasme
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
  -v, --verbose              verbose output
```

### SEE ALSO

* [wasme build](../wasme_build)	 - Build a wasm image from the filter source directory using Bazel-in-Docker
* [wasme catalog](../wasme_catalog)	 - interact with catalog
* [wasme deploy](../wasme_deploy)	 - Deploy an Envoy WASM Filter to the data plane (Envoy proxies).
* [wasme init](../wasme_init)	 - Initialize a project directory for a new Envoy WASM Filter.
* [wasme list](../wasme_list)	 - List Envoy WASM Filters stored locally or published to webassemblyhub.io.
* [wasme login](../wasme_login)	 - login so you can push images to webassemblyhub.io and submit them to the curated catalog
* [wasme pull](../wasme_pull)	 - Pull wasm filters from remote registry
* [wasme push](../wasme_push)	 - Push a wasm filter to remote registry
* [wasme tag](../wasme_tag)	 - Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
* [wasme undeploy](../wasme_undeploy)	 - Remove a deployed Envoy WASM Filter from the data plane (Envoy proxies).

