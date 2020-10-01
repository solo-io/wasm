
## WASM Artifact Image Specification v0.0.0

- [Introduction](#introduction)
- [Terminology](#terminology)
- [Description](#description)
    - [Overview](#overview)
    - [Layers](#layers)
    - [Running OCI Images with Envoy](#running-oci-images-with-envoy)
- [Format](#format)
- [Envoy WASM Filter Specification](#envoy-wasm-filter-specification)
    - [Example](#example)


### Introduction:

The WASM Artifact Image Specification defines how to bundle WASM modules as OCI images. WASM Artifact Images consist of a WASM binary file, configuration file, and metadata for the target WASM runtime.

The spec is intended to be generic, allowing for any type of WASM module whether it is used to extend any Envoy, OPA, or the browser.

The spec can be considered an extension of the OCI Artifact Spec designed specifically for use by applications which produce and consume WASM modules (as opposed to application containers). It is intended to provide a standard mechanism to manage the building and running of WASM modules. 

This document considers primarily the use case of storing WASM Envoy Filters as OCI Images.

### Terminology:

| Term                               | Definition                                       |
|------------------------------------|--------------------------------------------------|
| WASM Module                        | The distributable, loadable, and executable unit of code in WebAssembly. 
| WASM OCI Image Specification       | The specification for storing WASM modules as OCI Images.
| WASM Runtime                       | The execution environment into which a WASM Module may be loaded. This refers to the application itself which loads and executes a wasm module. Examples include web browsers, the Open Policy Agent, the Envoy Proxy, or any other application which supports extension via WASM modules. 
| Runtime Configuration              | Configuration specific to the runtime which consumes a module. This configuration is stored as JSON and bundled with the module in the image in the specification. 
| Envoy WASM Filter                  | Custom Filters for the Envoy Proxy built as a WASM module.
| Envoy WASM OCI Image               | Envoy Filters stored as OCI images according to the specification. 
| Envoy WASM OCI Artifact Specification | An extension of the WASM OCI Artifact Spec which describes how to bundle and ship Envoy WASM filters as OCI Images. |

### Description:

#### Overview:

The WASM OCI Artifact Specification defines a method of storing WASM modules which makes them easy to build, pull, publish, and execute.

Because each execution environment (runtime) for a WASM module may have runtime-specific configuration parameters, a WASM image is composed of both a content layer, for the WASM module itself, as well as a config layer, with metadata describing the module which is relevant to the target runtime.

#### Layers:

The content layer always consists of the WASM module binary. 

The config layer consists of a JSON-formatted string, which contains metadata for the target runtime. The runtime and ABI (Application Binary Interface) versions of an image can be deduced by parsing the config layer. 

The config layer may also contain additional data, depending on the target runtime. For example, the config for a WASM Envoy Filter contains root_ids available on the filter. 

For the sake of simplicity, the specification only supports a single module per image.

#### Running OCI Images with Envoy:

Envoy supports loading and running WASM modules via a file on local disk or an “Http datasource”.

Envoy WASM Filters can be stored according to the spec and run with Istio and Gloo, with the help of a local cache which pulls images from remote registries.  

Control planes then configure the Envoy instances to load the filter via the local cache, using the required root_id parameter supplied in the image config if it is available.


### Format:

The WASM OCI Artifact Spec consists of two layers bundled together:
- A layer specifying configuration for the target runtime
- A layer containing the compiled WASM module itself

Each layer is associated with its own Media Type, which is stored in the OCI Descriptor for that layer:

| Media Type | Type | Description |
|------------|------|-------------|
| application/vnd.io.solo.wasm.config.v1+json | JSON Object | Configuration for the Target WASM runtime.
| application/vnd.io.solo.wasm.content.layer.v1+wasm | binary data (byte array) | The compiled module data |

`application/vnd.io.solo.wasm.config.v1+json` Property Descriptions:

| Property   | Type | Description |
|------------|------|-------------|
| type | string | Name of the target runtime. Required. Specifies the intended runtime of the bundled module. The content of the Opaque JSON Config is specific to the type of WASM runtime. 
| abiVersions | string array | List of ABI Versions for the target runtime with which the image is compatible. The format for the version is dependent upon the runtime itself.
| config | JSON Object | This field stores any configuration parameters required by the target runtime. Its structure depends on the specified runtime. |


### Envoy WASM Filter Specification

The runtime config for Envoy WASM Filter OCI Images has the following format:

- *type* is set to `envoy_proxy`
- *abiVersion* is set to a recognized version of the Envoy Proxy WASM Filter ABI 
- *config* is a JSON Object containing a list of Filter root_ids that can be used with the provided filter

The `root_ids` key in the *config* JSON Object will have a list of strings as a value. Each string in the list corresponds to the name of a registered Root Context Helper defined in the module.

#### Example:

The following descriptors provide an example of the OCI Image descriptors for an Envoy WASM Filter stored according to the specification:
```
[
  {
    "mediaType": "application/vnd.io.solo.wasm.config.v1+json",
    "digest": "sha256:d0a165298ae270c5644be8e9938036a3a7a5191f6be03286c40874d761c18abf",
    "size": 125,
    "annotations": {
      "org.opencontainers.image.title": "runtime-config.json"
    }
  },
  {
    "mediaType": "application/vnd.io.solo.wasm.content.layer.v1+wasm",
    "digest": "sha256:5e82b945b59d03620fb360193753cbd08955e30a658dc51735a0fcbc2163d41c",
    "size": 1043056,
    "annotations": {
      "org.opencontainers.image.title": "filter.wasm"
    }
  }
]
```

The following is the runtime config stored as the `application/vnd.io.solo.wasm.config.v1+json` layer:

```{
  "type": "envoy_proxy",
  "abi_version": "v0-541b2c1155fffb15ccde92b8324f3e38f7339ba6",
  "config": {
    "root_ids": [
      "add_header_root_id"
    ]
  }
}
```

You can use the `wasme` tool to take new or existing module code and package it according to the WASM OCI Spec.