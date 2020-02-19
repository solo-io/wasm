---
title: "wasme build"
weight: 5
---
## wasme build

Build a wasm image from the filter source directory.

### Synopsis

Options for the build are specific to the target language.

### Options

```
  -c, --config string    The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.
  -h, --help             help for build
  -i, --image string     Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
      --store string     Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -t, --tag string       The image ref with which to tag this image. Specified in the format <name:tag>. Required
      --tmp-dir string   Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux. If unset, temporary files will be removed after build
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - The tool for building, pushing, and deploying Envoy WebAssembly Filters
* [wasme build assemblyscript](../wasme_build_assemblyscript)	 - Build a wasm image from an AssemblyScript filter using NPM-in-Docker
* [wasme build cpp](../wasme_build_cpp)	 - Build a wasm image from a CPP filter using Bazel-in-Docker
* [wasme build precompiled](../wasme_build_precompiled)	 - Build a wasm image from a Precompiled filter.

