
# Wasm Image Specification v0.0.0

## Introduction:

This document describes a varient of [Wasm Artifact Image Specification](README.md), which leverages the OCI "compatible" media type. Here, we omit definition and terminology explained in [Wasm Artifact Image Specification](README.md). 

We call this variant "compat", and the spec in [Wasm Artifact Image Specification](README.md) "oci".

## Description

This *compat* variant makes use of `application/vnd.oci.image.layer.v1.tar+gzip` media type for layers, and is not based on custom OCI Artifcat media types. This way users can oeperate with standard tools such as docker, podman, buildah, etc.

## Format

### Annotation

The *compat* variant must add the annotation `module.wasm.image/variant=compat` in the manifest.

### Layer

The *compat* variant must consist of exactly one `application/vnd.oci.image.layer.v1.tar+gzip` layer containing the two files:
- `plugin.wasm` - (**Required**) A Wasm binary to be loaded by the runtime.
- `runtime-config.json` - (**Optional**) A runtime configuratio specified in [Wasm Artifact Image Specification](README.md).

### Example

The following is an example OCI manifest of a *compat* variant image:

```
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:933594cea89247a78932eb9d74fae998e6fc3d1d114a7ff7705aaf702dbf7edb",
    "size": 326
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:e05c6f7d59f4c5976d9c1be8e12c34f64c49e5541967581e7f052070705191ac",
      "size": 151
    }
  ],
  "annotations": {
    "module.wasm.image/variant": "compat"
  }
}
```

And the contents in the layer consists of two files mentioned above

```
$ tar tf blobs/sha256/e05c6f7d59f4c5976d9c1be8e12c34f64c49e5541967581e7f052070705191ac
filter.wasm
runtime-config.json
```


## Appendix: build a *compat* image with Buildah

In this section, we demonstrate how to build a compiliant image with Buildah, a standard cli for building OCI images. We use v1.21.0 of Buildah here.

We assume that you have a valid Wasm binary named `filter.wasm` and `runtime-config.json` that you want to package as a Wasm OCI image.

First, we create a working container from `scratch` base image with `buildah from` command.

```
$ buildah --name mywasm from scratch
mywasm
```

Next, add the annotation described above via `buildah config` command

```
$ buildah config --annotation "module.wasm.image/variant=compat" mywasm
```

Then copy the files into that base image by `buildah copy` command to create the layer.

```
$ buildah copy mywasm filter.wasm runtime-config.json ./
af82a227630327c24026d7c6d3057c3d5478b14426b74c547df011ca5f23d271
```

Now, you can build a *compat* image and push it to your registries via `buildah commit` command

```
$ buildah commit mywasm docker://my-remote-registry/mywasm:0.1.0
```
