---
title: "wasme build precompiled"
weight: 5
---
## wasme build precompiled

Build a wasm image from a Precompiled filter.

### Synopsis


wasme supports building deployable images from a precompiled .wasm file. The user must provide their own configuration file with the --config flag.

The specification for this flag can be found here: [{{< versioned_link_path fromRoot="/reference/image_config">}}]({{< versioned_link_path fromRoot="/reference/image_config">}})


```
wasme build precompiled COMPILED_FILTER_FILE --tag <name:tag> --config <image config> [flags]
```

### Options

```
  -h, --help   help for precompiled
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

