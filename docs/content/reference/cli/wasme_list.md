---
title: "wasme list"
weight: 5
---
## wasme list

List Envoy WASM Filters stored locally or published to webassemblyhub.io.

### Synopsis

List Envoy WASM Filters stored locally or published to webassemblyhub.io.

```
wasme list [flags]
```

### Options

```
  -h, --help                            help for list
      --published                       Set to true to list images that have been published to a remote registry. If unset, lists images stored in local image cache.
      --search wasme list --published   Search images from the remote registry. If unset, wasme list --published will return the top repositories which are accessed the most.
  -s, --server string                   If using --published, read images from this remote registry. (default "webassemblyhub.io")
  -d, --show-dir                        Set to true to show the local directories for images. Does not apply to published images.
      --store string                    Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store. Ignored if using --published
  -w, --wide                            Set to true to list images with their full tag length.
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - The tool for building, pushing, and deploying Envoy WebAssembly Filters

