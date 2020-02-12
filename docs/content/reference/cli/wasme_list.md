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
  -h, --help            help for list
      --published       Set to true to list images that have been published to webassemblyhub.io. Defaults to listing image stored in local image cache.
  -s, --server string   If using --published, read images from this remote registry. (default "yuvaltest.solo.io")
      --store string    Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store. Ignored if using --published
  -w, --wide            Set to true to list images with their full tag length.
```

### Options inherited from parent commands

```
  -d, --debug     debug mode
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - 

