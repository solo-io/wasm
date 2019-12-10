---
title: "Getting Started with WebAssembly Hub"
weight: 1
description: Create a simple WebAssembly filter and run it with Envoy
---

In this tutorial we'll take a look at creating a new WebAssembly (WASM) module using the tooling from [WebAssembly Hub](https://webassemblyhub.io) to simplify a lot of the packaging and sharing aspects. We'll then load the new module into a running Envoy Proxy and verify it works. Lastly, we'll look at the workflow to get the new module into the WebAssembly Hub via a pull-request workflow. 

You can use your own distribution of Envoy that supports WebAssembly, but for this tutorial, we'll use Gloo which is an API Gateway based on Envoy. Gloo 1.0 announced beta support for WASM


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

```shell
$  kubectl apply -f gloo.yaml
```

*****
This section needs love... it's not working yet ^^ we need a nice way to install Gloo
****

### Verify set up

Lastly, we'll set up our routing rules to be able to call our `petstore` service. Let's add a route to the routing table:

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

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for getting the WebAssembly Hub CLI tools.

Let's create a new project called `new-filter`:

```shell
$  wasme init ./new-filter
```

We should now have a new folder with our new project:

```shell
$  cd new-filter
$  ls -l 
-rw-r--r--  1 ceposta  staff    572 Dec 10 09:06 BUILD
-rw-r--r--  1 ceposta  staff    505 Dec 10 09:06 README.md
-rw-r--r--  1 ceposta  staff   2782 Dec 10 09:06 WORKSPACE
drwxr-xr-x  3 ceposta  staff     96 Dec 10 09:06 bazel
-rw-r--r--  1 ceposta  staff   2797 Dec 10 09:06 filter.cc
-rw-r--r--  1 ceposta  staff     60 Dec 10 09:06 filter.proto
-rw-r--r--  1 ceposta  staff  30776 Dec 10 09:06 gloo-gateway-wasm.yaml
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

Now, let's build our filter. We can build using Bazel if we have that installed, or we can leverage a docker container that sets up all of the dependencies for us.

For Bazel directly:

```shell
$  bazel build :filter.wasm
```

From the docker container using the `wasme` tooling:

```shell
$  wasme build . 
```

This will download all of the necessary dependencies and compile the output `filter.wasm` into the `_output` folder if using the `wasme` tooling. Otherwise you can find filter in `./bazel-bin/filter.wasm`.

Now that we've built the `wasm` module, let's package it and load it into a registry so we can consume it in our Envoy proxy.

## Pushing wasm module to registry

Before we use our new filter.wasm module, let's push it into a registry that can be used as a source for the module when we configure Envoy. Alternatively we could try load it from the disk, but that involves a lot more.

To do that, let's login to the `webassemblyhub.io` using GitHub as the OAuth provider. From the CLI:

```shell
$  wasm login

Using port: 60632
Opening browser for login. If the browser did not open for you, please go to:  https://webassemblyhub.io/authorize?port=60632
```

You should see a GitHub OAuth screen:

![](/img/wasme_login.png)

Click the "Authorize" button at the bottom and continue.

After successful auth, you should see this in the terminal:

```shell
success ! you are now authenticated
```

Now let's push to the webassemblyhub.io registry. 

```shell
$  wasme push webassemblyhub.io/christian-posta/test:v0.1 ./_output_/filter.wasm
```

NOTE: The tag name to use is
`webassemblyhub.io/<your-git-username>/<whatever-name>:<whatever-version>`

When you've pushed, you should be able to see your new module in the registry:

```shell
$  wasme list

NAME                        SHA      UPDATED             SIZE   TAGS
christian-posta/test        6aef37f3 13 Jan 10 12:54 MST 1.0 MB v0.1
```

## Deploy our new module

To deploy this new WASM module into Envoy, using Gloo we can configure the gateway to load the module from the `webassembly.io` registry:

```yaml
apiVersion: gateway.solo.io.v2/v2
kind: Gateway
metadata:
  labels:
    app: gloo
  name: gateway-proxy-v2
  namespace: gloo-system
spec:
  bindAddress: '::'
  bindPort: 8080
  httpGateway:
    plugins:
      extensions:
        configs:
          wasm:
            image: webassemblyhub.io/christian-posta/test:v0.1
            name: christian
            root_id: add_header_root_id
  proxyNames:
  - gateway-proxy-v2
  useProxyProto: false
```

If we `kubectl apply -f ` this configuration to our cluster, we should expect the Envoy proxy to pick up this new module and dynamically load it into the proxy. 

Note, it could take a few seconds for the module to get picked up.

## Verify behavior

To verify the behavior, let's use the `curl` command from above:

```shell
curl -v $URL/api/pets
```

We expect to see our new headers in the response:

```shell
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
< date: Tue, 10 Dec 2019 16:43:19 GMT
< content-length: 86
< x-envoy-upstream-service-time: 2
< newheader: 
< doc-header: it-worked
< location: envoy-wasm
< server: envoy
< 
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
```
You can see our new header `doc-header` in our response!

## Add this new module to the Web Assembly Hub catalog

Just like when we pushed the module to the registry, to add the module to the catalog on WebAssembly Hub, we need to be logged in. If you're already logged in, skip this command:

```shell
$  wasme login
```

Adding the the catalog is a git-ops workflow that involves a Pull Request. To initiate the process, run the `wasme catalog add` command:


```shell
$  wasme catalog add webassemblyhub.io/christian-posta/test:v0.1
```

Note, the tag we are adding matches what we built in previous steps.

This will give us an interactive prompt to fill in the details of our WebAssembly Hub extension:

```shell
? Please provide name of the extension test
? Please provide short description of the extension (required) short des
? Please provide long description of the extension long desc
? Please provide the url to the source code https://gitub.com/christian-posta/test
? Please provide the url to the documentation https://url.com
? Please provide the name of the extension creator christian-posta
? Please provide a url for the extension creator https://blog.christianposta.com
? Please provide a url for the extension logo http://logo.com
```

Lastly, you'll be prompted one last time whether to proceed. Type `Y` and hit enter:

```shell
creatorName: christian-posta
creatorUrl: https://blog.christianposta.com
documentationUrl: https://url.com
extensionRef: webassemblyhub.io/christian-posta/test:v0.1
logoUrl: http://logo.com
longDescription: long desc
name: test
repositoryUrl: https://gitub.com/christian-posta/test
shortDescription: short des

In these steps we will:
         Fork github.com/solo-io/wasme
         Create a feature branch named test-v0.1
         Add your spec to this branch in this location: catalog/test/v0.1/spec.yaml
         And open a Pull Request against github.com/solo-io/wasme
 Yes
Making sure your fork is available
Making sure a feature branch is available
Adding your catalog item to the feature branch
Makeing sure a PR is open
Created PR: https://github.com/solo-io/wasme/pull/17
```

If you open the `https://github.com/solo-io/wasme/pull/17` URL, you should see a PR was created in the `wasme` GitHub repo:

![](/img/test-pr.png)

## Pull Request

The Pull Request from here will be reviewed by the WebAssembly Hub team and included in the catalog if it's accepted. 
