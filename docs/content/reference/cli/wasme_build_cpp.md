---
title: "wasme build cpp"
weight: 5
---
## wasme build cpp

Build a wasm image from a CPP filter using Bazel-in-Docker

### Synopsis

Build a wasm image from a CPP filter using Bazel-in-Docker

```
wasme build cpp SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag> [flags]
```

### Options

```
  -f, --bazel-ouptut bazel-bin   Path relative to bazel-bin to the wasm file produced by running the Bazel target. (default "filter.wasm")
  -g, --bazel-target string      Name of the bazel target to run. (default ":filter.wasm")
  -b, --build-dir string         Directory containing the target BUILD file. (default ".")
  -h, --help                     help for cpp
```

### Options inherited from parent commands

```
  -c, --config string    The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.
  -d, --debug            debug mode
  -i, --image string     Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
      --store string     Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -t, --tag string       The image ref with which to tag this image. Specified in the format <name:tag>. Required
      --tmp-dir string   Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux. If unset, temporary files will be removed after build
  -v, --verbose          verbose output
```

### SEE ALSO

* [wasme build](../wasme_build)	 - Build a wasm image from the filter source directory.

