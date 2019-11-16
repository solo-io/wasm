This repo contains the tooling that allows you to build and push WASM envoy filters,
so they will be accessible to gloo.

To do so, first build the tool
```
go build .
```

The build the example filter (mostly copied from the envoy-wasm):
```
(cd example; bazel build :envoy_filter_http_wasm_example.wasm --config=wasm)
```

Push:
```
./extend-envoy push gcr.io/solo-public/yuvalism:v1  example/bazel-bin/yuval
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
            image: gcr.io/solo-public/yuvalism:v1
            name: yuval
            root_id: my_root_id
```