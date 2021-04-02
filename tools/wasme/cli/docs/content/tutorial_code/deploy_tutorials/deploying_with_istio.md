---
title: "Deploying Wasm Filters to Istio"
weight: 2
description: Deploy a wasm filter using Istio as the control plane.
---

Using Envoy's Web Assembly capabilities, we can add custom filters to an [Istio](https://istio.io) service mesh. This allows us to customize and extend functionality of the mesh's [data plane](https://blog.envoyproxy.io/service-mesh-data-plane-vs-control-plane-2774e720f7fc).

In this tutorial we'll use `wasme` to deploy a simple "hello world" filter that adds a header to HTTP responses. This WebAssembly (WASM) module has already been built and can be pulled from [the WebAssembly Hub](https://webassemblyhub.io). To learn how to build and push filters, [see the tutorial on building and pushing wasm filters](../getting_started.md).


## Prepare environment

### Install Istio

1. First, we'll download the latest Istio release. We take version 1.8.1 as an example:

```bash
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.8.1 sh -
cd istio-1.8.1
```

1. To install Istio:

```bash
istioctl install --set profile=demo
```

You should get output like this:

```bash
This will install the Istio 1.8.1 demo profile with ["Istio core" "Istiod" "Ingress gateways" "Egress gateways"] components into the cluster. Proceed? (y/N) y
✔ Istio core installed                                                                                                                
✔ Istiod installed                                                                                                                    
✔ Egress gateways installed                                                                                                           
✔ Ingress gateways installed                                                                                                          
✔ Installation complete 
```

### Deploy Bookinfo App

1. First, we'll deploy the Istio [`bookinfo` Application](https://istio.io/docs/examples/bookinfo/):

```bash
kubectl create ns bookinfo
kubectl label namespace bookinfo istio-injection=enabled --overwrite
kubectl apply -n bookinfo -f samples/bookinfo/platform/kube/bookinfo.yaml 
```

### Testing the Setup

To ensure everything installed correctly, let's try running a request between two of the deployed services:

```bash
# execute a request from the productpage component to the details component: 
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

{{% notice note %}}
It may take a few minutes before all the Istio sidecars are ready to serve traffic.
{{% /notice %}}

The output should look have a `200 OK` response and look like the following:

{{< highlight yaml "hl_lines=9-15" >}}
*   Trying 10.102.48.118...
* TCP_NODELAY set
* Connected to details.bookinfo (10.102.48.118) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Fri, 02 Apr 2021 07:23:31 GMT
< content-length: 180
< x-envoy-upstream-service-time: 32
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

You should get output like this:

```
NAME                                      TAG  SIZE    SHA      UPDATED
webassemblyhub.io/ilackarms/add-header v0.1 12.6 kB 0295d929 02 Apr 21 13:06 CST
```

Deploying the filter is done with a single `wasme` command:

```bash
wasme deploy istio webassemblyhub.io/ilackarms/add-header:v0.1 \
    --id=myfilter \
    --namespace bookinfo \
    --config 'world'
```

{{% notice note %}}
The `config` for the `webassemblyhub.io/ilackarms/add-header:v0.1` filter specifies the value of 
a "`hello`" header which will be appended by the filter to HTTP responses. The value of `config` is specific to the 
filter deployed via `wasme`.
{{% /notice %}}


Wasme will output the following logs as it deploys the filter to each `deployment` that composes the bookinfo app:

```
INFO[0000] cache namespace already exists                cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0000] cache configmap already exists                cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0000] cache service account already exists          cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0000] cache role updated                            cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0000] cache rolebinding updated                     cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0000] cache daemonset updated                       cache=wasme-cache.wasme image="quay.io/solo-io/wasme:0.0.33"
INFO[0007] image is already cached                       cache="{wasme-cache wasme}" image="webassemblyhub.io/tanjunchen20/add-header:v0.1"
INFO[0007] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=details-v1
INFO[0007] created Istio EnvoyFilter resource            envoy_filter_resource=details-v1-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=details-v1
INFO[0008] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=productpage-v1
INFO[0008] created Istio EnvoyFilter resource            envoy_filter_resource=productpage-v1-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=productpage-v1
INFO[0008] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=ratings-v1
INFO[0008] created Istio EnvoyFilter resource            envoy_filter_resource=ratings-v1-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=ratings-v1
INFO[0009] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v1
INFO[0009] created Istio EnvoyFilter resource            envoy_filter_resource=reviews-v1-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v1
INFO[0009] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v2
INFO[0009] created Istio EnvoyFilter resource            envoy_filter_resource=reviews-v2-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v2
INFO[0009] updated workload sidecar annotations          filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v3
INFO[0009] created Istio EnvoyFilter resource            envoy_filter_resource=reviews-v3-myfilter.bookinfo filter="id:\"myfilter\" image:\"webassemblyhub.io/tanjunchen20/add-header:v0.1\" config:<type_url:\"type.googleapis.com/google.protobuf.StringValue\" value:\"\\n\\005world\" > rootID:\"add_header\" patchContext:\"inbound\" " workload=reviews-v3
```

If the above command finished without error, we should be ready to test the filter:

```bash
# execute a request from the productpage component to the details component: 
kubectl exec -ti -n bookinfo deploy/productpage-v1 -c istio-proxy -- curl -v http://details.bookinfo:9080/details/123
```

The output should look have a `200 OK` response and contain the response header `hello: world`:

{{< highlight yaml "hl_lines=15" >}}
*   Trying 10.108.142.139...
* TCP_NODELAY set
* Connected to details.bookinfo (10.108.142.139) port 9080 (#0)
> GET /details/123 HTTP/1.1
> Host: details.bookinfo:9080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-type: application/json
< server: istio-envoy
< date: Fri, 02 Apr 2021 09:44:48 GMT
< content-length: 180
< x-envoy-upstream-service-time: 3
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
