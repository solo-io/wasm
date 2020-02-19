---
title: "wasme pull"
weight: 5
---
## wasme pull

Pull wasm filters from remote registry

### Synopsis

Pull wasm filters from remote registry


```
wasme pull <name:tag|name@digest> [flags]
```

### Options

```
  -c, --config stringArray   path to auth config
  -h, --help                 help for pull
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

* [wasme](../wasme)	 - 

