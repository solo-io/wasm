---
title: "wasme build tinygo"
weight: 5
---
## wasme build tinygo

Build a wasm image from a TinyGo filter using TinyGo-in-Docker

### Synopsis

Build a wasm image from a TinyGo filter using TinyGo-in-Docker

```
wasme build tinygo SOURCE_DIRECTORY -t <name:tag> [flags]
```

### Options

```
  -h, --help   help for tinygo
```

### Options inherited from parent commands

```
  -c, --config string    The path to the filter configuration file for the image. If not specified, defaults to <SOURCE_DIRECTOR>/runtime-config.json. This file must be present in order to build the image.
  -i, --image string     Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
      --store string     Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -t, --tag string       The image ref with which to tag this image. Specified in the format <name:tag>. Required
      --tmp-dir string   Directory for storing temporary files during build. Defaults to /tmp on OSx and Linux. If unset, temporary files will be removed after build
  -v, --verbose          verbose output
```

### SEE ALSO

* [wasme build](../wasme_build)	 - Build a wasm image from the filter source directory.

