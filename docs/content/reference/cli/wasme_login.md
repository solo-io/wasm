---
title: "wasme login"
weight: 5
---
## wasme login

Log in so you can push images to the remote server.

### Synopsis


Caches credentials for image pushes in the provided credentials-file (defaults to $HOME/.wasme/credentials.json).

Provide -s=SERVER_ADDRESS to provide login credentials for a registry other than webassemblyhub.io.



```
wasme login [-s SERVER_ADDRESS] -u USERNAME -p PASSWORD  [flags]
```

### Options

```
      --credentials-file string   write to this credentials file. defaults to $HOME/.wasme/credentials.json
  -h, --help                      help for login
  -p, --password string           login password
      --plaintext                 use plaintext to connect to the remote registry (HTTP) rather than HTTPS
  -s, --server string             the address of the remote registry to which to authenticate (default "webassemblyhub.io")
  -u, --username string           login username
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [wasme](../wasme)	 - The tool for building, pushing, and deploying Envoy WebAssembly Filters

