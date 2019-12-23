---
title: "Deploying Filters to Gloo"
weight: 2
description: Deploy a wasm filter using Gloo as the control plane.
---

In this tutorial we'll deploy an existing WebAssembly (WASM) module from [the WebAssembly Hub](https://webassemblyhub.io) directly to Envoy via [Gloo](https://docs.solo.io/gloo/latest) installed to our kubernetes cluster.

## Prepare environment

To get started, let's deploy a sample service that we can call through Envoy. We'll deploy the sample petstore API:

```shell
kubectl apply -f \
https://raw.githubusercontent.com/solo-io/gloo/master/example/petstore/petstore.yaml
```

You should now have the petstore running:

```shell
kubectl get po 
```

```

NAME                        READY   STATUS    RESTARTS   AGE
petstore-5dcf5d6b66-n8tjt   1/1     Running   0          2m20s
```

### Deploying Envoy

In this tutorial, we'll use Gloo, an API Gateway based on Envoy that has built-in wasm support but these steps should also work for base Envoy.

First, install Gloo via the helm chart:

```shell
helm repo update
kubectl create ns gloo-system
helm install gloo-gateway gloo/gloo --namespace gloo-system \
  --set global.wasm.enabled=true
```

Gloo will be installed to the `gloo-system` namespace.

### Verify set up

Lastly, we'll set up our routing rules to be able to call our `petstore` service. Let's add a route to the routing table:

Download and apply the [virtual service manifest](../default-virtualservice.yaml)
```shell
kubectl apply -f default-virtualservice.yaml
```

To get Gloo's external IP, run the following:

```shell
URL=$(kubectl get svc -n gloo-system gateway-proxy \
 -o jsonpath='{.status.loadBalancer.ingress[*].ip}')
```

Now let's curl that URL:

```shell
    curl -v $URL/api/pets
```

```

*   Trying 35.184.102.75...
* TCP_NODELAY set
* Connected to 35.184.102.75 (35.184.102.75) port 80 (#0)
> GET /api/pets HTTP/1.1
> Host: 35.184.102.75
> User-Agent: curl/7.54.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< content-type: application/xml
< date: Tue, 10 Dec 2019 16:02:00 GMT
< content-length: 86
< x-envoy-upstream-service-time: 2
< server: envoy
< 
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```

If you're able to get to this point, we have a working Envoy proxy and we're able to call it externally. 

## Deploying a WASM module from the Hub

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for getting the WebAssembly Hub CLI `wasme`.

Let's run `wasme list` to see what's available on the hub:

```shell
wasme list
```

```
NAME                 SHA      UPDATED             SIZE   TAGS
...
ilackarms/hello:v0.1 3753eeaf 15 Sep 19 23:41 EST 1.0 MB v0.1
...
```

Let's try deploying one of these to Gloo:

```bash
wasme deploy gloo webassemblyhub.io/ilackarms/hello:v0.1 --id=myfilter
```

This filter adds the header `hello: World!` to responses.

The deployment should have added our filter to the Gloo Gateway. Let's check this with `kubectl`:

```bash
kubectl get gateway -n gloo-system '-ojsonpath={.items[0].spec.httpGateway.options.wasm}'
```

```
map[image:webassemblyhub.io/ilackarms/hello:v0.1 name:myfilter rootId:add_header_root_id]
```

If we try our request again, we should see the `hello: World` header was added by our filter:


```shell
curl -v $URL/api/pets
```

```
*   Trying 34.73.225.160...
* TCP_NODELAY set
* Connected to 34.73.225.160 (34.73.225.160) port 80 (#0)
> GET /api/pets HTTP/1.1
> Host: 34.73.225.160
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< content-type: application/xml
< date: Fri, 20 Dec 2019 19:17:13 GMT
< content-length: 86
< x-envoy-upstream-service-time: 0
< hello: World!
< location: envoy-wasm
< server: envoy
<
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```

Note that deploying filters is dynamic and does not require restarting the proxy. 

## Cleaning up

We can clean up our filter with the `wasme undeploy` command:

```bash
wasme undeploy gloo --id=myfilter
```

Then re-try the `curl`:

```shell
curl -v $URL/api/pets
```

```
*   Trying 34.73.225.160...
* TCP_NODELAY set
* Connected to 34.73.225.160 (34.73.225.160) port 80 (#0)
> GET /api/pets HTTP/1.1
> Host: 34.73.225.160
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< content-type: application/xml
< date: Fri, 20 Dec 2019 19:19:13 GMT
< content-length: 86
< x-envoy-upstream-service-time: 1
< server: envoy
<
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
* Connection #0 to host 34.73.225.160 left intact
```

Cool! We've just seen how easy it is to dynamically add and remove filters from Envoy using `wasme`.

For more information and support using `wasme` and the Web Assembly Hub, visit the Solo.io slack channel at
https://slack.solo.io.
