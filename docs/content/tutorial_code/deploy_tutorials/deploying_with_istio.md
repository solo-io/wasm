---
title: "Deploying Wasm Filters to Istio"
weight: 2
description: Deploy a wasm filter using Istio as the control plane.
---

Using Envoy's Web Assembly capabilities, we can add custom filters to an [Istio](https://istio.io) service mesh. This allows us to customize and extend functionality of the mesh's [data plane](https://blog.envoyproxy.io/service-mesh-data-plane-vs-control-plane-2774e720f7fc).

In this tutorial we'll use `wasme` to deploy a simple "hello world" filter that adds a header to HTTP responses. This WebAssembly (WASM) module has already been built and can be pulled from [the WebAssembly Hub](https://webassemblyhub.io). To learn how to build and push filters, [see the tutorial on building and pushing wasm filters](../getting_started.md).


## Prepare environment

### Install Istio

1. First, we'll download the latest Istio release. At time of writing, this is `1.4.2`:

```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.4.2
```

1. To install Istio:

```bash
bin/istioctl manifest apply --set profile=demo
```  

### Deploy Bookinfo App

1. First, we'll deploy the Istio [`bookinfo` Application](https://istio.io/docs/examples/bookinfo/):

```bash
kubectl create ns bookinfo
kubectl label namespace bookinfo istio-injection=enabled --overwrite
kubectl apply -n bookinfo \
  -f samples/bookinfo/platform/kube/bookinfo.yaml 
```

### Testing the Setup

To ensure everything installed correctly, let's try running a request between two of the deployed services:

```bash
# execute a request from the productpage component to the details component: 
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

The output should look have a `200 OK` response and look like the following:

{{< highlight yaml "hl_lines=9-15" >}}
*   Trying 10.55.247.3...
* TCP_NODELAY set
* Connected to details.bookinfo (10.55.247.3) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Mon, 06 Jan 2020 16:28:55 GMT
< content-length: 180
< x-envoy-upstream-service-time: 2
< x-envoy-decorator-operation: details.bookinfo.svc.cluster.local:9080/*
<
* Connection #0 to host details.bookinfo left intact
{"id":123,"author":"William Shakespeare","year":1595,"type":"paperback","pages":200,"publisher":"PublisherA","language":"English","ISBN-10":"1234567890","ISBN-13":"123-1234567890"}
{{< /highlight >}}

In the next section, we'll add a simple filter to the bookinfo sidecars.  

## Deploy the filter

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for getting the WebAssembly Hub CLI `wasme`.

Let's run `wasme list` to see what's available on the hub:

```shell
wasme list
```

```
NAME                          SHA      UPDATED             SIZE   TAGS
...
ilackarms/istio-example       9d9be851 02 May 85 11:02 EST 1.0 MB 1.4.2
...
```

Deploying the filter is done with a single `wasme` command:

```bash
wasme deploy istio webassemblyhub.io/ilackarms/istio-test:1.4.2-0 \
    --id=myfilter \
    --namespace bookinfo \
    --config '{"name":"hello","value":"world"}'
```

{{% notice note %}}
The `config` for the `webassemblyhub.io/ilackarms/istio-example` filter specifies the `name` and `value` of 
a header which will be appended by the filter to HTTP responses. The value of `config` is specific to the 
filter deployed via `wasme`. 
{{% /notice %}}


Wasme will output the following logs as it deploys the filter to each `deployment` that composes the bookinfo app:

```
INFO[0000] cache namespace already exists                cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.1"
INFO[0000] cache configmap already exists                cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.1"
INFO[0000] cache daemonset updated                       cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.1"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0001] added image to cache                          cache="{wasme-cache wasme}"
INFO[0002] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0002] updated Istio EnvoyFilter resource            envoy_filter_resource=details-v1-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
INFO[0002] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0003] updated Istio EnvoyFilter resource            envoy_filter_resource=productpage-v1-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
INFO[0003] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0003] updated Istio EnvoyFilter resource            envoy_filter_resource=ratings-v1-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
INFO[0003] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0004] updated Istio EnvoyFilter resource            envoy_filter_resource=reviews-v1-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
INFO[0004] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0004] updated Istio EnvoyFilter resource            envoy_filter_resource=reviews-v2-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
INFO[0004] updated workload sidecar annotations          filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
webassemblyhub.io/ilackarms/istio-test:1.4.2-0 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:eec1bcf3c5fdf1d5e4e69d923ebbc7c5e06f409d6f3f914aac3b386895e45eac Size:14 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:e6e9b6296a2cb633fb7c8206ebfb2771752684180065de03a86aaa591ebeee1f Size:1039813 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:9d9be851b8c937583ffa78cceb902541b039e5c2138c8d391d426394a18f61b0 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0005] updated Istio EnvoyFilter resource            envoy_filter_resource=reviews-v3-myfilter.bookinfo filter="&{myfilter webassemblyhub.io/ilackarms/istio-test:1.4.2-0 {\"name\":\"hello\",\"value\":\"world\"} }" ports="[9080]" workload="{ bookinfo deployment}"
```

If the above command finished without error, we should be ready to test the filter:


```bash
# execute a request from the productpage component to the details component: 
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

The output should look have a `200 OK` response and contain the response header `hello: world`:

{{< highlight yaml "hl_lines=15" >}}
*   Trying 10.55.247.3...
* TCP_NODELAY set
* Connected to details.bookinfo (10.55.247.3) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Mon, 06 Jan 2020 18:13:12 GMT
< content-length: 180
< x-envoy-upstream-service-time: 1
< hello: world
< x-envoy-decorator-operation: details.bookinfo.svc.cluster.local:9080/*
<
* Connection #0 to host details.bookinfo left intact
{"id":123,"author":"William Shakespeare","year":1595,"type":"paperback","pages":200,"publisher":"PublisherA","language":"English","ISBN-10":"1234567890","ISBN-13":"123-1234567890"}
{{< /highlight >}}

Cool! We've just deployed a Web Assembly filter to Envoy with a single command!
 
To remove the filter, run: 

```bash 
wasme undeploy istio --id myfilter --namespace bookinfo
```

For more information and support using `wasme` and the Web Assembly Hub, visit the Solo.io slack channel at
https://slack.solo.io.
