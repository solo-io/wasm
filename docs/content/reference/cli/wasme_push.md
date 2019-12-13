---
title: "wasme push"
weight: 5
---
## wasme push

Push wasm filter to remote registry

### Synopsis

Push wasm filter to remote registry. E.g.:

wasme push webassemblyhub.io/my/filter:v1 filter.wasm


```
wasme push name[:tag|@digest] code.wasm [config_proto-descriptor-set.proto.bin] [flags]
```

### Options

```
  -d, --debug     debug mode
  -h, --help      help for push
  -v, --verbose   verbose output
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

