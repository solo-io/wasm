---
title: "wasme build"
weight: 5
---
## wasme build

Build a wasm image from the filter source directory using Bazel-in-Docker

### Synopsis

Build a wasm image from the filter source directory using Bazel-in-Docker

```
wasme build SOURCE_DIRECTORY [-b <bazel target>] [-t <name:tag>] [flags]
```

### Options

```
  -f, --bazel-ouptut bazel-bin   Path relative to bazel-bin to the wasm file produced by running the Bazel target. (default "filter.wasm")
  -g, --bazel-target string      Name of the bazel target to run. (default ":filter.wasm")
  -b, --build-dir string         Directory containing the target BUILD file. (default ".")
  -h, --help                     help for build
  -i, --image string             Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
      --store string             Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -t, --tag string               The image ref with which to tag this image. Specified in the format <name:tag>. Required
      --wasm-file string         If specified, wasme will use the provided path to a compiled filter wasm to produce the image. The bazel build will be skipped and the wasm-file will be used instead.
```

### Options inherited from parent commands

```
  -c, --config stringArray   auth config path
  -d, --debug                debug mode
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
  -v, --verbose              verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - 

