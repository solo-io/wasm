---
title: "wasme build assemblyscript"
weight: 5
---
## wasme build assemblyscript

Build a wasm image from an AssemblyScript filter using NPM-in-Docker

### Synopsis

Build a wasm image from an AssemblyScript filter using NPM-in-Docker

```
wasme build assemblyscript SOURCE_DIRECTORY [-b <bazel target>] -t <name:tag> [flags]
```

### Options

```
  -e, --email string      Email for logging in to NPM before running npm install. Optional
  -h, --help              help for assemblyscript
  -p, --password string   Password for logging in to NPM before running npm install. Optional
  -u, --username string   Username for logging in to NPM before running npm install. Optional
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

