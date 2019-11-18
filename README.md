This repo contains the tooling that allows you to build and push WASM envoy filters,
so they will be accessible to gloo.

To do so, first build the tool
```
go build .
```

The build the example filter (mostly copied from the envoy-wasm):
```
(cd example; bazel build :envoy_filter_http_wasm_example.wasm)
```

Push:
```
./extend-envoy push gcr.io/solo-public/example-filter:v1 example/bazel-bin/envoy_filter_http_wasm_example.wasm
```

load in to gloo:
```
kubectl edit -n gloo-system gateways.gateway.solo.io.v2 gateway-proxy-v2
```

set the httpGateway field like so:
```
  httpGateway:
    plugins:
      extensions:
        configs:
          wasm:
            config: yuval
            image: gcr.io/solo-public/example-filter:v1
            name: yuval
            root_id: my_root_id
```


# emscripten sdk
If you change the emscripten SDK, an sdk with PR merged is needed:
https://github.com/emscripten-core/emscripten/pull/9812/files