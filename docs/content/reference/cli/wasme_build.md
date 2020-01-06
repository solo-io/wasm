---
title: "wasme build"
weight: 5
---
## wasme build

Compile the filter to wasm using Bazel-in-Docker

### Synopsis

Compile the filter to wasm using Bazel-in-Docker

```
wasme build SOURCE_DIRECTORY [-b <bazel target>] [-o OUTPUT_FILE] [flags]
```

### Options

```
  -f, --bazel-ouptut bazel-bin   Path relative to bazel-bin to the wasm file produced by running the Bazel target. (default "filter.wasm")
  -t, --bazel-target string      Name of the bazel target to run. (default ":filter.wasm")
  -b, --build-dir string         Directory containing the target BUILD file. (default ".")
  -h, --help                     help for build
  -i, --image string             Name of the docker image containing the Bazel run instructions. Modify to run a custom builder image (default "quay.io/solo-io/ee-builder:dev")
  -o, --output string            path to the output .wasm file. Nonexistent directories will be created. (default "filter.wasm")
```

### Options inherited from parent commands

```
  -c, --config stringArray   auth config path
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
  -u, --username string      registry username
```

### SEE ALSO

* [wasme](../wasme)	 - 

