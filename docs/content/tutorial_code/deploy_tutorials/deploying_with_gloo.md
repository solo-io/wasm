---
title: "Deploying Filters to Gloo"
weight: 2
description: Deploy a wasm filter to Envoy using Gloo as the control plane.
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

{{% notice note %}}
Gloo version `1.3.6` or greater required. Check your installed version of Gloo with `glooctl version`
{{% /notice %}}

### Create a Route

Lastly, we'll create a route to be able to call our `petstore` service. Let's add a route using a Gloo `VirtualService`:

Apply the [virtual service manifest](https://docs.solo.io/gloo/latest/gloo_routing/virtual_services/)
```shell
cat <<EOF | kubectl apply -f -
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
      - prefix: /
      routeAction:
        single:
          upstream:
            name: default-petstore-8080
            namespace: gloo-system
EOF
```

Next, we'll get the Gloo Gateway's external IP by running the following:

```shell
URL=$(kubectl get svc -n gloo-system gateway-proxy \ -o jsonpath='{.status.loadBalancer.ingress[*].ip}')
```

{{< tabs >}}
{{< tab name="Cloud Provider" codelang="shell">}}
URL=$(kubectl get svc -n gloo-system gateway-proxy \ -o jsonpath='{.status.loadBalancer.ingress[*].ip}')
{{< /tab >}}
{{< tab name="Minikube" codelang="shell" >}}
URL=$(minikube ip):$(kubectl get svc -n gloo-system gateway-proxy -o jsonpath='{.spec.ports[?(@.name == "http")].nodePort}')`
{{< /tab >}}
{{< /tabs >}}

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

If you're able to get to this point, we are able to call to our petstore app through Gloo.

## Deploying a WASM module from the Hub

If you don't have `wasme` installed, try the following or refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for getting the WebAssembly Hub CLI `wasme`:

```bash
curl -sL https://run.solo.io/wasme/install | sh
export PATH=$HOME/.wasme/bin:$PATH
```

Let's run `wasme list` to see what's available on the hub:

```shell
wasme list --published
```

```
NAME                                   TAG                                 SIZE    SHA      UPDATED
...
webassemblyhub.io/ilackarms/gloo-test  1.3.3-0                             1.0 MB  8c001279 12 Feb 20 19:10 UTC
...
```

Let's try deploying one of these to Gloo:

```bash
wasme deploy gloo webassemblyhub.io/ilackarms/gloo-test:1.3.3-0 --id=myfilter --config 'world'
```

This filter adds the header `hello: <value>` to responses, where `<value>` is the value of the `--config` string.

The deployment should have added our filter to the Gloo Gateway. Let's check this with `kubectl`:

```bash
kubectl get gateway -n gloo-system '-ojsonpath={.items[0].spec.httpGateway.options.wasm}'
```

```
map[config:world image:webassemblyhub.io/ilackarms/gloo-test:1.3.3-0 name:myfilter rootId:add_header_root_id]
```

If we try our request again, we should see the `hello: world` header was added by our filter:


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
< hello: world
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
