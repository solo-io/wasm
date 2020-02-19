---
title: "Building WASM Filters in AssemblyScript"
weight: 1
description: "Build a simple WebAssembly filter in AssemblyScript."
---

In this tutorial we will write an Envoy filter in [AssemblyScript](https://docs.assemblyscript.org/) and build it using `wasme`.
 
 We'll optionally push
the image to the public WASM registry at https://webassemblyhub.io/.

## Creating a new WASM module

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for installing `wasme`, the WebAssembly Hub CLI.

Let's create a new project called `new-filter`:

```shell
$  wasme init ./assemblyscript-filter
```

You'll be asked with an interactive prompt which language platform you are building for. Choose the appropriate option below:

{{< tabs >}}
{{< tab name="istio" codelang="shell">}}
? What language do you wish to use for the filter:
  ▸ assemblyscript
? With which platform do you wish to use the filter?:
  ▸ istio 1.5.x
{{< /tab >}}
{{< tab name="gloo" codelang="shell" >}}
? What language do you wish to use for the filter:
  ▸ assemblyscript
? With which platform do you wish to use the filter?:
  ▸ gloo 1.3.x
{{< /tab >}}
{{< /tabs >}}

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
wasme build . -t webassemblyhub.io/ilackarms/add-header:v0.1
```

The module will take up to a few minutes to build. In the background, `wasme` has launched a Docker container to run the necessary 
build steps. 

When the build has finished, you'll be able to see the image with `wasme list`:

```bash
wasme list
```

```
NAME                                     SHA      UPDATED             SIZE   TAGS
webassemblyhub.io/ilackarms/add-header  bbfdf674 26 Jan 20 10:45 EST 1.0 MB v0.1
```

## Next Steps

Now that we've successfully built our image, we can try [running it locally]({{< versioned_link_path fromRoot="/tutorial_code/deploy_tutorials/deploying_with_local_envoy">}}) or [pushing it to a remote registry]({{< versioned_link_path fromRoot="/tutorial_code/push_tutorials">}}) so it can be pulled and deployed in a Kubernetes environment.
