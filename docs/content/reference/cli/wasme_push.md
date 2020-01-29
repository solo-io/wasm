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
  -h, --help           help for push
      --store string   Set the path to the local storage directory for wasm images. Defaults to $HOME/.wasme/store
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

