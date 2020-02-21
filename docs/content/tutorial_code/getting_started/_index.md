---
title: "Getting Started"
weight: 1
description: "Create a simple WebAssembly filter and deploy it to Envoy."
---

In this tutorial we will:

1. Write our own custom filter for Envoy
2. `build` a WASM module from our filter and store it as an OCI image/
3. `push` the module to the [WebAssembly Hub](https://webassemblyhub.io)
4. `deploy` the image to a running instance of Envoy.
4. `curl` the instance to see our filter act on a request.

For in-depth guides, please refer to:

- [Build tutorials](../build_tutorials) for topics relating to *building* WASM filters 
- [Deployment tutorials](../deploy_tutorials) for topics relating to *deploying* WASM filters 

## Prepare environment

To get started, let's deploy a sample service that we can call through Envoy. We'll deploy the sample petstore API:

```shell
$  kubectl apply -f \
https://raw.githubusercontent.com/solo-io/gloo/master/example/petstore/petstore.yaml
```

You should now have the petstore running:

```shell
$  kubectl get po 

NAME                        READY   STATUS    RESTARTS   AGE
petstore-5dcf5d6b66-n8tjt   1/1     Running   0          2m20s
```

### Deploying Envoy

In this tutorial, we'll use Gloo, an API Gateway based on Envoy that has built-in wasm support but these steps should also work for base Envoy.

First, install Gloo using one of the following installation options: 

{{< tabs >}}
{{< tab name="install-gloo" codelang="shell">}}
helm repo add gloo https://storage.googleapis.com/solo-public-helm
helm repo update
kubectl create ns gloo-system
helm install --namespace gloo-system --set global.wasm.enabled=true gloo gloo/gloo
{{< /tab >}}
{{< tab name="glooctl" codelang="shell" >}}
glooctl install gateway -n gloo-system --values <(echo '{"namespace":{"create":true},"crds":{"create":true},"global":{"wasm":{"enabled":true}}}')
{{< /tab >}}
{{< /tabs >}}

{{% notice note %}}
You can install `glooctl` via 
```
curl -sL https://run.solo.io/gloo/install | sh
export PATH=$HOME/.gloo/bin:$PATH
```
{{% /notice %}}

Gloo will be installed to the `gloo-system` namespace.

{{% notice note %}}
You can deploy your own gloo (version 1.2.10 and above), by enabling the experimental WASM support when 
installing. When installing you need to set the "global.wasm.enabled" flag to true. If installing
with glooctl, you can use the following command:
```shell
glooctl install gateway -n gloo-system --values <(echo '{"namespace":{"create":true},"crds":{"create":true},"global":{"wasm":{"enabled":true}}}')
```
You can add the `--dry-run` flag to glooctl to generate a yaml for you instead of installing directly.
{{% /notice %}}

### Verify set up

Lastly, we'll set up our routing rules to be able to call our `petstore` service. Let's add a route to the routing table:

Download and apply the [virtual service manifest](default-virtualservice.yaml)
```shell
$  kubectl apply -f default-virtualservice.yaml
```

To get Gloo's external IP, run the following:

```shell
$  URL=$(kubectl get svc -n gloo-system gateway-proxy \
 -o jsonpath='{.status.loadBalancer.ingress[*].ip}')
```

Now let's curl that URL:

```shell
$  curl -v $URL/api/pets

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

## Creating a new WASM module

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for installing `wasme`, the WebAssembly Hub CLI.

Let's create a new project called `new-filter`:

```shell
$  wasme init ./new-filter
```

You'll be asked with an interactive prompt which language platform you are building for.

Select `cpp` and `gloo-1.3.x`:

```noop
? What language do you wish to use for the filter:
  ▸ cpp
? With which platform do you wish to use the filter?:
  ▸ gloo 1.3.x
```

We should now have a new folder with our new project:

```shell
$  cd new-filter
$  ls -l 
-rw-r--r--  1 ceposta  staff    572 Dec 10 09:06 BUILD
-rw-r--r--  1 ceposta  staff    505 Dec 10 09:06 README.md
-rw-r--r--  1 ceposta  staff   2782 Dec 10 09:06 WORKSPACE
drwxr-xr-x  3 ceposta  staff     96 Dec 10 09:06 bazel
drwxr-xr-x  3 ceposta  staff     34 Dec 10 09:06 filter-config.json
-rw-r--r--  1 ceposta  staff   2797 Dec 10 09:06 filter.cc
-rw-r--r--  1 ceposta  staff     60 Dec 10 09:06 filter.proto
drwxr-xr-x  7 ceposta  staff    224 Dec 10 09:06 toolchain
```

Open this project in your favorite IDE. The source code is C++ and we'll make some changes to create a new filter.


## Making changes to the sample filter

Feel free to explore the project that was created. This is a Bazel project that sets up the correct emscripten toolchain to build C++ into WASM. Let's open ./new-filter/filter.cc in our favorite IDE.

Navigate to the `FilterHeadersStatus AddHeaderContext::onResponseHeaders` member method. Let's add a new header that we can use to verify our module was executed correctly (later down in the tutorial). Let's add a new response header named `doc-header`:

```code
addResponseHeader("doc-header", "it-worked");
```
Your method should look like this:

```go
FilterHeadersStatus AddHeaderContext::onResponseHeaders(uint32_t) {
  LOG_DEBUG(std::string("onResponseHeaders ") + std::to_string(id()));
  addResponseHeader("newheader", root_->header_value_);
  addResponseHeader("doc-header", "it-worked");
  replaceResponseHeader("location", "envoy-wasm");
  return FilterHeadersStatus::Continue;
}
```

## Building the filter

Now, let's build a WASM image from our filter with `wasme`. The filter will be tagged and stored
in a local registry, similar to how [Docker](https://www.docker.com/) stores images. 

In this example we'll include the registry address `webassemblyhub.io` as well as 
our GitHub username which will be used to authenticate to the registry.

Build and tag our image like so:

```shell
wasme build . -t ilackarms/add-header:gloo-1.3
```

The module will take up to a few minutes to build. In the background, `wasme` has launched a Docker container to run the necessary 
build steps. 

When the build has finished, you'll be able to see the image with `wasme list`:

```bash
wasme list
```

```
NAME                                    SHA      UPDATED             SIZE   TAGS
webassemblyhub.io/ilackarms/add-header  bbfdf674 26 Jan 20 10:45 EST 1.0 MB gloo-1.3
```

Now that we've built the WASM module, let's publish it into a registry so we can deploy it to our Envoy proxy running in Kubernetes.

## Pushing WASM module to registry

To do that, let's login to the `webassemblyhub.io` using GitHub as the OAuth provider. From the CLI:

```shell
$  wasme login

Using port: 60632
Opening browser for login. If the browser did not open for you, please go to:  https://webassemblyhub.io/authorize?port=60632
```

You should see a GitHub OAuth screen:

![](../../img/wasme_login.png)

Click the "Authorize" button at the bottom and continue.

After successful auth, you should see this in the terminal:

```shell
success ! you are now authenticated
```

Now let's push to the webassemblyhub.io registry. 

```shell
$  wasme push webassemblyhub.io/ilackarms/add-header:gloo-1.3
INFO[0000] Pushing image webassemblyhub.io/ilackarms/add-header:gloo-1.3
INFO[0001] Pushed webassemblyhub.io/ilackarms/add-header:gloo-1.3
INFO[0001] Digest: sha256:22f2d81f9b61ebbf1aaeb00aa7bde20a90b90ec8bb5f492cc18a140de205dc32
```

{{% notice note %}}
The tag name to use is
`webassemblyhub.io/<your-git-username>/<some-name>:<some-version>`
{{% /notice %}}

When you've pushed, you should be able to see your new module in the registry:

```shell
$  wasme list --published  

NAME                        SHA      UPDATED             SIZE   TAGS
ilackarms/add-header        6aef37f3 13 Jan 10 12:54 MST 1.0 MB gloo-1.3
```

## Deploy our new module

To deploy the module to Envoy via Gloo:

```bash
wasme deploy gloo webassemblyhub.io/ilackarms/add-header:gloo-1.3 --id=add-header \
  --config '{"name":"hello","value":"World!"}'
```

It will take a few moments for the image to be pulled by the server-side cache.

## Verify behavior

To verify the behavior, let's use the `curl` command from above:

```shell
curl -v $URL/api/pets
```

We expect to see our new headers in the response:

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

You can see our new header `hello: World!` in our response!
