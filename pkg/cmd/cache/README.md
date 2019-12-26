Deploying cache for usage with Istio:

# Building a filter

Istio 1.4.2 is using a proxy based on https://github.com/istio/envoy/ commit 6d525c67f39b36cdff9d688697f266c1b55e9cb7 
(sha 53efbe4a35e2c4dbee9a2ee50ce6469eb8bbcbb997bec5a2023623d60f629c0c)

To build a filter, you can start off from the example in the `example/cpp-istio-1.4` folder.

build with:
```
bazel build :filter.wasm
```

push it to somewhere publicly accessible. webassemblyhub.io for example
```
wasme push webassemblyhub.io/yuval-k/istio-example:1.4.2 bazel-bin/filter.wasm
```

# Build wasme cache container
**Note** this is only needed if you changed the cache code.
If you update the cache code, build and push docker image:
```
docker build -t quay.io/solo-io/wasme:0.0.1 .
docker push quay.io/solo-io/wasme:0.0.1
```

# Deploy cache container

```
kubectl create namespace wasme-cache
kubectl apply -n wasme-cache -f pkg/cmd/cache/deploy.yaml
```

# Edit deployment where filter is desired

Edit your pod to include the cached filters as a volume:

```
sidecar.istio.io/userVolume: '[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache"}}]'
sidecar.istio.io/userVolumeMount: '[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]'
```

Note that these should be on the annotations of the pod template: e.g.:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: details-v1
  labels:
    app: details
    version: v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: details
      version: v1
  template:
    metadata:
      annotations:
        sidecar.istio.io/userVolume: '[{"name":"cache-dir","hostPath":{"path":"/var/local/lib/wasme-cache","type":"Directory"}}]'
        sidecar.istio.io/userVolumeMount: '[{"mountPath":"/var/local/lib/wasme-cache","name":"cache-dir"}]'
      labels:
        app: details
        version: v1
    spec:
      serviceAccountName: bookinfo-details
      containers:
      - name: details
        image: docker.io/istio/examples-bookinfo-details-v1:1.15.0
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 9080
```

you can use the following command to apply:
```
kubectl -n default patch deploy/details-v1 --type=merge -p='{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/userVolume":"[{\"name\":\"cache-dir\",\"hostPath\":{\"path\":\"/var/local/lib/wasme-cache\",\"type\":\"Directory\"}}]","sidecar.istio.io/userVolumeMount":"[{\"mountPath\":\"/var/local/lib/wasme-cache\",\"name\":\"cache-dir\"}]"}}}}}'
```

# Apply filter - EnvoyFilter CRD
Add an envoy filter CRD with your filter. By default the deployment has a test wasme filter. it is cached with this `bbfdf674f5cf2e413a4b701ae865dd4569502f60af5647f2f47ca4f38e2b40af` hash.

Note: the `wasme deploy` command will soon be updated to automate this step of providing the filter hash.

```
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: details-wasm
spec:
  workloadSelector:
    labels:
      app: details
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      listener:
        portNumber: 9080
        filterChain:
          filter:
            name: "envoy.http_connection_manager"
            subFilter:
              name: "envoy.router"
    patch:
      operation: INSERT_BEFORE
      value:
        name: envoy.filters.http.wasm
        config:
          config:
            name: "test"
            root_id: "add_header_root_id"
            configuration: '{"value":"hi-from-wasme"}'
            vm_config:
              runtime: envoy.wasm.runtime.v8
              code:
                filename:
                  /var/local/lib/wasme-cache/bbfdf674f5cf2e413a4b701ae865dd4569502f60af5647f2f47ca4f38e2b40af
```

apply with:

```
kubectl apply -n default -f pkg/cmd/cache/filter.yaml
```

Test:

```
kubectl proxy &
curl -v http://localhost:8001/api/v1/namespaces/default/services/details:9080/proxy/details/123
```

# debugging

check details config dump:
```
kubectl port-forward -n default deploy/details-v1 15000 &
curl http://localhost:15000/config_dump
```

logs from envoy:
```
curl 'localhost:15000/logging?level=debug' -XPOST
kubectl logs -n default deploy/details-v1 -c istio-proxy
```

logs from pilot (see nacks in there):
```
kubectl logs -n istio-system deploy/istio-pilot -c discovery
```