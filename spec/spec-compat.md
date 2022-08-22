
# Wasm Image Specification v0.0.0

## Introduction:

This document describes a variant of [Wasm Artifact Image Specification](spec.md), which leverages the compatible layer media types. Here, we omit definition and terminology explained in [Wasm Artifact Image Specification](spec.md). 

We call this variant "compat", and the spec in [Wasm Artifact Image Specification](spec.md) "oci".

## Description

This *compat* variant makes use of compatible media type for layers, and is not based on custom OCI Artifact media types. This way users can operate with standard tools such as docker, podman, buildah, and standard container registries which don't yet support custom media types as used in [*oci* variant](spec.md).

## Specification

### Layer

The *compat* variant must have as last layer one whose media type is one of the followings:
- `application/vnd.oci.image.layer.v1.tar+gzip`
- `application/vnd.docker.image.rootfs.diff.tar.gzip`

In addition, the layer must consist of the following two files:
- `plugin.wasm` - (**Required**) A Wasm binary to be loaded by the runtime.
- `runtime-config.json` - (**Optional**) A runtime configuration specified in [Wasm Artifact Image Specification](spec.md#Format). This is mainly used as metadata for the image, depending on the runtime.

### Annotation

If the media type equals `application/vnd.oci.image.layer.v1.tar+gzip`, then a *compat* variant image *should* add the annotation `module.wasm.image/variant=compat` in the manifest to make it easy to distinguish this *compat* variant from the *oci* variant. Note that this is **optional**.

### Example with `application/vnd.oci.image.layer.v1.tar+gzip` media type

The following is an example OCI manifest of images with `application/vnd.oci.image.layer.v1.tar+gzip` layer media type:

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

### Example with `application/vnd.docker.image.rootfs.diff.tar.gzip` media type

The following is an example Docker manifest of images with `application/vnd.docker.image.rootfs.diff.tar.gzip` layer media type:

```
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 1182,
    "digest": "sha256:500c5c9b0755790c440f6d24a8926e399913bda2d599dcac24edb99a72b66de7"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 161,
      "digest": "sha256:cf72304d01ead8fe014ed9f09e4132678ee4f29030ec46e6242c457071435ec3"
    }
  ]
}
```

## Appendix 1: build a *compat* image with Buildah

We demonstrate how to build a *compat* image with Buildah, a standard cli for building OCI images. We use v1.21.0 of Buildah here. Produced images have `application/vnd.oci.image.layer.v1.tar+gzip` layer media type

We assume that you have a valid Wasm binary named `plugin.wasm` and `runtime-config.json` (optional) that you want to package as an image.

1. First, we create a working container from `scratch` base image with `buildah from` command.

```
$ buildah --name mywasm from scratch
mywasm
```

2. Next, add the annotation described above via `buildah config` command

```
$ buildah config --annotation "module.wasm.image/variant=compat" mywasm
```

**Note this step is optional. See [Annotation](#annotation) section.**


3. Then copy the files into that base image by `buildah copy` command to create the layer.

```
$ buildah copy mywasm plugin.wasm runtime-config.json ./
af82a227630327c24026d7c6d3057c3d5478b14426b74c547df011ca5f23d271
```

**Note: you must execute `buildah copy` exactly once in order to end up having only one layer in produced images**

4. Now, you can build a *compat* image and push it to your registry via `buildah commit` command

```
$ buildah commit mywasm docker://my-remote-registry/mywasm:0.1.0
```

## Appendix 2: build a *compat* image with Docker CLI

We demonstrate how to build a *compat* image with Docker CLI. Produced images have `application/vnd.docker.image.rootfs.diff.tar.gzip` layer media type.

We assume that you have a valid Wasm binary named `plugin.wasm` and `runtime-config.json` (optional) that you want to package as an image.

1. First, we prepare the following Dockerfile:

```
$ cat Dockerfile
FROM scratch
LABEL org.opencontainers.image.title my-wasm-extension
COPY runtime-config.json plugin.wasm ./
```

**Note: you must have exactly one `COPY` instruction in the Dockerfile at the end as only the last layer in produced images is going to be taken into account to obtain the files**

2. Then, build your image via `docker build` command

```
$ docker build . -t my-registry/mywasm:0.1.0
```

3. Finally, push the image to your registry via `docker push` command

```
$ docker push my-registry/mywasm:0.1.0
```
