
---
title: "module.wasm.configruntime.proto"
---

## Package : `module.wasm.config`



<a name="top"></a>

<a name="API Reference for runtime.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## runtime.proto


## Table of Contents
  - [EnvoyConfig](#module.wasm.config.EnvoyConfig)
  - [Runtime](#module.wasm.config.Runtime)







<a name="module.wasm.config.EnvoyConfig"></a>

### EnvoyConfig
configuration for an Envoy Filter WASM Image


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| root_ids | [][string](#string) | repeated | the set of root IDs exposed by the Envoy Filter |






<a name="module.wasm.config.Runtime"></a>

### Runtime
Runtime Configuration for a WASM OCI Image. This configuration is bundled
with the WASM image at build time.

Example:

```json
{
  &#34;type&#34;: &#34;envoy_proxy&#34;,
  &#34;abiVersions&#34;: [&#34;v0-541b2c1155fffb15ccde92b8324f3e38f7339ba6&#34;],
  &#34;config&#34;: {
    &#34;rootIds&#34;: [
      &#34;add_header_root_id&#34;
    ]
  }
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | the type of the runtime |
| abi_versions | [][string](#string) | repeated | the compatible versions of the ABI of the target runtime
this may be different than the version of the runtime itself
this is used to ensure compatibility with the runtime |
| config | [EnvoyConfig](#module.wasm.config.EnvoyConfig) |  | the config for running the module
currently, wasme only supports Envoy config |





 

 

 

 

