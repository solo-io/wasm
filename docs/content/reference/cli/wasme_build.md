---
title: "wasme build"
weight: 5
---
## wasme build

Compile the filter to wasm using Bazel-in-Docker

### Synopsis

Compile the filter to wasm using Bazel-in-Docker

```
wasme build SOURCE_DIRECTORY [-o OUTPUT_FILE] [flags]
```

### Options

```
  -h, --help            help for build
  -i, --image string    Name of the docker image containing the Bazel run instructions. Only be modified if you are an experiejnced user (default "quay.io/solo-io/ee-builder:v1")
  -o, --output string   path to the output .wasm file. Nonexistent directories will be created. (default "_output/filter.wasm")
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

