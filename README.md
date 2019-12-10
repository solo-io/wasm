This repo contains the tooling that allows you to build and push WASM envoy filters,
so they will be accessible to gloo.

To do so, first build the tool
```
go build .
```

The build the example filter (mostly copied from the envoy-wasm):
```
(cd example; bazel build :filter.wasm :filter_proto)
```

Push:
```
./wasme push gcr.io/solo-public/example-filter:v1 example/bazel-bin/filter.wasm example/bazel-bin/filter_proto-descriptor-set.proto.bin
```

load in to gloo:
```
kubectl edit -n gloo-system gateways.gateway.solo.io gateway-proxy
```

set the httpGateway field like so:
```
  httpGateway:
    options:
      wasm:
        config: |
          {}
        image: webassemblyhub.io/yuval-k/metrics:v1
        name: yuval
        root_id: stats_root_id
```

Download the petstore from the following tutorial https://docs.solo.io/gloo/latest/gloo_routing/hello_world/

Then apply the following virtual service to enable routing to the petstore.

```
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default
  namespace: gloo-system
spec:
  virtualHost:
    domains:
    - '*'
    routes:
    - matchers:
        - exact: /sample-route-1
      routeAction:
        single:
          upstream:
            name: default-petstore-8080
            namespace: gloo-system
      options:
        prefixRewrite: /api/pets
```

Now call the API with the following command
```bash
$ curl $(glooctl proxy url)/sample-route-1


[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```

Congrats, you officially used a WASM filter.
This is a simple stats filter so all it is doing is updating some basic prometheus statistics on the routes living on this listener.
For more complex and interesting filters check out https://webassemblyhub.io/.

# emscripten sdk
If you change the emscripten SDK, an sdk with PR merged is needed:
https://github.com/emscripten-core/emscripten/pull/9812/files
