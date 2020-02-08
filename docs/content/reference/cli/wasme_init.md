---
title: "wasme init"
weight: 5
---
## wasme init

Initialize a project directory for a new Envoy WASM Filter.

### Synopsis


The provided --language flag will determine the programming language used for the new filter. The default is 
C++.

The provided --platform flag will determine the target platform used for the new filter. This is important to 
ensure compatibility between the filter and the 

If --language, --platform, or --platform-version are not provided, the CLI will present an interactive prompt. Disable the prompt with --disable-prompt



```
wasme init DEST_DIRECTORY [--language=FILTER_LANGUAGE] [--platform=TARGET_PLATFORM] [--platform-version=TARGET_PLATFORM_VERSION] [flags]
```

### Options

```
      --disable-prompt            Disable the interactive prompt if a required parameter is not passed. If set to true and one of the required flags is not provided, wasme CLI will return an error.
  -h, --help                      help for init
      --language string           The programming language with which to create the filter. Supported languages are: [cpp]
      --platform string           The name of the target platform against which to build. Supported platforms are: [gloo istio]
      --platform-version string   The version of the target platform against which to build. Supported Istio versions are: [1.4.x 1.5.x]. Supported Gloo versions are: [1.3.x]
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

