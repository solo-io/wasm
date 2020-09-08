---
title: "wasme push"
weight: 5
---
## wasme push

Push a wasm filter to remote registry

### Synopsis

Push wasm filter to remote registry. E.g.:

wasme push webassemblyhub.io/my/filter:v1


```
wasme push name[:tag|@digest] [flags]
```

### Options

```
  -c, --config stringArray   path to auth config
  -h, --help                 help for push
      --insecure             allow connections to SSL registry without certs
  -p, --password string      registry password
      --plain-http           use plain http and not https
      --store string         Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
  -u, --username string      registry username
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - The tool for building, pushing, and deploying Envoy WebAssembly Filters

