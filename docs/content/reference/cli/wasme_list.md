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
  -h, --help           help for list
      --published      Set to true to list images that have been published to webassemblyhub.io. Defaults to listing image stored in local image cache.
      --store string   Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store. Ignored if using --published
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

