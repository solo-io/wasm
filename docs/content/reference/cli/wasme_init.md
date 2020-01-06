---
title: "wasme init"
weight: 5
---
## wasme init

Initialize a source directory for new Envoy WASM Filter.

The provided --base will determine the content of the created directory. The default is 
a C++ example filter compatible with the latest Envoy Wasm APIs.

Note that Istio 1.4 uses an older version of the Envoy Wasm APIs and users should 
use --base=cpp-istio to initialize a filter source directory for Istio.


### Synopsis

Initialize a source directory for new Envoy WASM Filter.

The provided --base will determine the content of the created directory. The default is 
a C++ example filter compatible with the latest Envoy Wasm APIs.

Note that Istio 1.4 uses an older version of the Envoy Wasm APIs and users should 
use --base=cpp-istio to initialize a filter source directory for Istio.


```
wasme init DEST_DIRECTORY [--base=FILTER_BASE] [flags]
```

### Options

```
      ----base string   The type of filter to build. Valid filter bases are: [cpp cpp-istio] (default "cpp")
  -h, --help            help for init
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

