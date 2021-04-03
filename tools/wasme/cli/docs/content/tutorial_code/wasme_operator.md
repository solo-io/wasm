---
title: "Running Filters in Production"
weight: 6
description: For CI/CD workflows and production deployments of web assembly filters, the Wasme Operator offers a declarative custom resource for managing deployed wasm filters.
---

The `wasme` CLI provides an easy way to get started building and deploying Web Assembly filters to an Envoy service mesh.

This is intended to be used in development and testing, but does not provide a declarative, stateless means by which to configure production Kubernetes clusters.

The **Wasme Operator** makes it possible to manage the deployment of WebAssembly Filters to a supported service mesh using Kubernetes CRDs.

The Wasme Operator consists of two components:

- an *image cache*, which
    * pulls and caches wasm filter images from a compatible filter registry (such as `https://webassemblyhub.io/`)
    * is deployed as a Kubernetes DaemonSet (to make images available on all nodes)
- and an [*operator*](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/), which 
    * installs and configures wasm filters to the data plane proxies    
    * is deployed as a Kubernetes Deployment

All components run in the `wasme` namespace by default. 

# Installation

### Installing Istio

Wasme depends on a supported service mesh being installed to the cluster.

Currently, the Wasme Operator only supports Istio. If Istio (specifically the `envoyfilters.networking.istio.io` CRD) is not installed to the cluster, the wasme container will exit with an error at boot time.

To install istio, follow the [Istio installation Guide](https://istio.io/docs/setup/getting-started/#install).

{{% notice note %}}
Note that this guide was written and tested against Istio 1.8.2. Some older versions of Istio (<1.7.x) may have a different `config` type to configure the filter, which previously took a `string` value instead of a `google.protobuf.Any`.

If you are running this guide against a different minor version of Istio, it is recommended that you use `wasme init` to generate your own wasm filter targeting your specific Istio version.
{{% /notice %}}

### Installing Wasme

First, install the Wasme CRDs:

```bash
kubectl apply -f https://github.com/solo-io/wasm/releases/latest/download/wasme.io_v1_crds.yaml
```

Output:

```
customresourcedefinition.apiextensions.k8s.io/filterdeployments.wasme.io created
```

Next install the Operator components:

```bash
kubectl apply -f https://github.com/solo-io/wasm/releases/latest/download/wasme-default.yaml
```

Output:

```
namespace/wasme created
configmap/wasme-cache created
serviceaccount/wasme-cache created
serviceaccount/wasme-operator created
clusterrole.rbac.authorization.k8s.io/wasme-operator created
clusterrolebinding.rbac.authorization.k8s.io/wasme-operator created
daemonset.apps/wasme-cache created
deployment.apps/wasme-operator created
```

{{% notice note %}}
To install an older version of wasme, use the url `kubectl apply -f https://github.com/solo-io/wasm/releases/download/<VERSION>/wasme-default.yaml`
{{% /notice %}}

Finally, confirm that the wasme operator is has started successfully:

```bash
kubectl get pod -n wasme
```

Output:

```
NAME                              READY   STATUS    RESTARTS   AGE
wasme-cache-5twpj                 1/1     Running   0          4m40s
wasme-operator-754bb5f654-5wd6h   1/1     Running   0          4m40s
```

Great! We're now ready to get started deploying WebAssembly filters to our Istio service mesh!

See the next section to learn how to get started with the Operator.

# Using the Wasme Operator 

Interacting with the Wasme Operator happens through the `FilterDeployment`  Custom Resource.

The full spec for this CRD can be read at https://github.com/solo-io/wasm/blob/master/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto#L13

Let's try the following example to see it in action:

## Example Usage

In this example, we'll deploy a simple application with Istio sidecars injected to it. We'll deploy a simple "Hello World" filter to our application's sidecars and see that it modifies request headers accordingly.

#### Deploy the Example

For our example we'll use the [Istio Bookinfo example](https://istio.io/docs/examples/bookinfo/).

To deploy it, let's run the following:

```bash
# create the bookinfo namespace
kubectl create ns bookinfo

# label it for istio injection
kubectl label namespace bookinfo istio-injection=enabled --overwrite

# install the bookinfo application
kubectl apply -n bookinfo -f https://raw.githubusercontent.com/solo-io/wasm/master/tools/wasme/cli/test/e2e/operator/bookinfo.yaml
```

{{% notice note %}}
The bookinfo app installed here is identical to that shipped with Istio `1.8.2`.
{{% /notice %}}

#### Deploy the Filter

Deploying our filter to the Bookinfo sidecars is as simple as creating a **FilterDeployment** custom resource.

Let's take a brief look at an example FilterDeployment:

```yaml
apiVersion: wasme.io/v1
kind: FilterDeployment
metadata:
  name: bookinfo-custom-filter
  namespace: bookinfo
spec:
  deployment:
    istio:
      kind: Deployment
  filter:
    config:
      '@type': type.googleapis.com/google.protobuf.StringValue
      value: world
    image: webassemblyhub.io/sodman/istio-1-7:v0.3
```

This resource tells wasme to:

- add the `webassemblyhub.io/sodman/istio-1-7:v0.3` filter to each **Deployment** in the `bookinfo` namespace
- with the *configuration* string `world`


Run the following to add the filter to the Bookinfo app:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: wasme.io/v1
kind: FilterDeployment
metadata:
  name: bookinfo-custom-filter
  namespace: bookinfo
spec:
  deployment:
    istio:
      kind: Deployment
  filter:
    config:
      '@type': type.googleapis.com/google.protobuf.StringValue
      value: world
    image: webassemblyhub.io/sodman/istio-1-7:v0.3
EOF
```

```
filterdeployment.wasme.io/bookinfo-custom-filter created
```

The Wasme Operator will immediately begin processing the FilterDeployment. We should see its `status` is updated within a few seconds:

```bash
kubectl get filterdeployments.wasme.io -n bookinfo -o yaml bookinfo-custom-filter 
```

Note the `status` of the FilterDeployment

{{< highlight yaml "hl_lines=20-26" >}}
apiVersion: wasme.io/v1
kind: FilterDeployment
metadata:
  creationTimestamp: "2021-04-03T14:11:05Z"
  generation: 1
  name: bookinfo-custom-filter
  namespace: bookinfo
  resourceVersion: "7767415"
  selfLink: /apis/wasme.io/v1/namespaces/bookinfo/filterdeployments/bookinfo-custom-filter
  uid: 64f08103-3492-4e1b-a463-668e398204f6
spec:
  deployment:
    istio:
      kind: Deployment
  filter:
    config:
      '@type': type.googleapis.com/google.protobuf.StringValue
      value: world
    image: webassemblyhub.io/sodman/istio-1-7:v0.3
status:
  observedGeneration: "1"
  workloads:
    details-v1:
      state: Succeeded
    productpage-v1:
      state: Succeeded
{{< /highlight >}}

The `status` contains the status of the deployment for each selected workload. This means the filter has been deployed to each.

We check whether the corresponding envoyfilter is generated under the bookinfo namespace.

Let's get the envoyfilter in bookinfo namespace: 

```bash
kubectl get envoyfilter -n bookinfo
```

```
NAME                                             AGE
details-v1-bookinfo-custom-filter.bookinfo       6m13s
productpage-v1-bookinfo-custom-filter.bookinfo   6m12s
```

Let's test the filter with a `curl`:

```bash
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

The output should have a `200 OK` response and contain the response header `hello: world`:

{{< highlight yaml "hl_lines=15" >}}
*   Trying 10.107.216.139...
* TCP_NODELAY set
* Connected to details.bookinfo (10.107.216.139) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Tue, 09 Feb 2021 21:41:30 GMT
< content-length: 180
< x-envoy-upstream-service-time: 2
< hello: world
< location: envoy-wasm
< x-envoy-decorator-operation: details.bookinfo.svc.cluster.local:9080/*
< 
* Connection #0 to host details.bookinfo left intact
{"id":123,"author":"William Shakespeare","year":1595,"type":"paperback","pages":200,"publisher":"PublisherA","language":"English","ISBN-10":"1234567890","ISBN-13":"123-1234567890"}
{{< /highlight >}}

We can easily modify the `hello: world` custom header by updating the FilterDeployment `spec.filter.config.value`:

{{< highlight yaml "hl_lines=12" >}}
cat <<EOF | kubectl apply -f -
apiVersion: wasme.io/v1
kind: FilterDeployment
metadata:
  name: bookinfo-custom-filter
  namespace: bookinfo
spec:
  deployment:
    istio:
      kind: Deployment
  filter:
    config:
      '@type': type.googleapis.com/google.protobuf.StringValue
      value: goodbye
    image: webassemblyhub.io/sodman/istio-1-7:v0.3
EOF
{{< /highlight >}}

You can check that the envoyfilter `productpage-v1-bookinfo-custom-filter.bookinfo` in bookinfo namespace was updated with the new configuration value:

{{< highlight yaml "hl_lines=40" >}}
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  creationTimestamp: "2021-04-03T14:12:16Z"
  generation: 2
  name: details-v1-bookinfo-custom-filter.bookinfo
  namespace: bookinfo
  ownerReferences:
  - apiVersion: wasme.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: FilterDeployment
    name: bookinfo-custom-filter
    uid: 64f08103-3492-4e1b-a463-668e398204f6
  resourceVersion: "7768328"
  selfLink: /apis/networking.istio.io/v1alpha3/namespaces/bookinfo/envoyfilters/details-v1-bookinfo-custom-filter.bookinfo
  uid: c749b91c-02d6-491c-8314-b7f6faf02f8b
spec:
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_INBOUND
      listener:
        filterChain:
          filter:
            name: envoy.http_connection_manager
            subFilter:
              name: envoy.router
    patch:
      operation: INSERT_BEFORE
      value:
        name: envoy.filters.http.wasm
        typedConfig:
          '@type': type.googleapis.com/udpa.type.v1.TypedStruct
          typeUrl: type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm
          value:
            config:
              configuration:
                '@type': type.googleapis.com/google.protobuf.StringValue
                value: goodbye
              name: bookinfo-custom-filter.bookinfo
              rootId: add_header_root_id
              vmConfig:
                code:
                  local:
                    filename: /var/local/lib/wasme-cache/d2bc5bea58499684981fda875101ac18a69923cea4a4153958dad08065aa1e74
                runtime: envoy.wasm.runtime.v8
                vmId: bookinfo-custom-filter.bookinfo
  workloadSelector:
    labels:
      app: details
      version: v1
{{< /highlight >}}

Try the request again:

```bash
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

The output should now contain the response header `hello: goodbye`:

{{< highlight yaml "hl_lines=15" >}}
*   Trying 10.107.216.139...
* TCP_NODELAY set
* Connected to details.bookinfo (10.107.216.139) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Tue, 09 Feb 2021 21:45:42 GMT
< content-length: 180
< x-envoy-upstream-service-time: 2
< hello: goodbye
< location: envoy-wasm
< x-envoy-decorator-operation: details.bookinfo.svc.cluster.local:9080/*
< 
* Connection #0 to host details.bookinfo left intact
{"id":123,"author":"William Shakespeare","year":1595,"type":"paperback","pages":200,"publisher":"PublisherA","language":"English","ISBN-10":"1234567890","ISBN-13":"123-1234567890"}
{{< /highlight >}}

Great! We've just seen how easy it is to deploy Wasm filters to Istio using Wasme!

To remove the filter, run: 

```bash 
kubectl delete filterdeployment -n bookinfo bookinfo-custom-filter
```

For more information and support using `wasme` and the Web Assembly Hub, visit the Solo.io slack channel at
https://slack.solo.io.
