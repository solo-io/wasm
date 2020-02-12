---
title: "wasme build precompiled"
weight: 5
---
## wasme build precompiled

Build a wasm image from a Precompiled filter.

### Synopsis

Build a wasm image from a Precompiled filter.

```
wasme build precompiled COMPILED_FILTER_FILE -t <name:tag> [flags]
```

### Options

```
  -h, --help   help for precompiled
```

### Options inherited from parent commands

```
  -c, --config string      The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.
  -d, --debug              debug mode
  -i, --image string       Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
      --store string       Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -t, --tag string         The image ref with which to tag this image. Specified in the format <name:tag>. Required
      --tmp-dir string     Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux. If unset, temporary files will be removed after build
  -v, --verbose            verbose output
      --wasm-file string   If specified, wasme will use the provided path to a compiled filter wasm to produce the image. The bazel build will be skipped and the wasm-file will be used instead.
```

### SEE ALSO

* [wasme build](../wasme_build)	 - Build a wasm image from the filter source directory.

