---
title: "wasme pull"
weight: 5
---
## wasme pull

Pull files from remote registry

### Synopsis

Pull files from remote registry

Example - Pull only files with the "application/vnd.oci.image.layer.v1.tar" media type (default):
  oras pull localhost:5000/hello:latest

Example - Pull only files with the custom "application/vnd.me.hi" media type:
  oras pull localhost:5000/hello:latest -t application/vnd.me.hi

Example - Pull all files, any media type:
  oras pull localhost:5000/hello:latest -a

Example - Pull files from the insecure registry:
  oras pull localhost:5000/hello:latest --insecure

Example - Pull files from the HTTP registry:
  oras pull localhost:5000/hello:latest --plain-http


```
wasme pull <name:tag|name@digest> [-o output-file] [flags]
```

### Options

```
  -d, --debug           debug mode
  -h, --help            help for pull
  -o, --output string   output file (default "filter.wasm")
  -v, --verbose         verbose output
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

