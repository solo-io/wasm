---
title: "Building WASM Filters in C++"
weight: 1
description: "Build a simple C++ filter in AssemblyScript."
---

In this tutorial we will write an Envoy filter in C++ and build it using `wasme`.

## Creating a C++ WASM module

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for installing `wasme`, the WebAssembly Hub CLI.

Let's create a new filter called `cpp-filter`:

```shell
wasme init cpp-filter
```

You'll be asked with an interactive prompt which language platform you are building for. At time of writing, `wasme` includes separate bases 
 for Istio 1.5.x and Gloo 1.3.x:

{{< tabs >}}
{{< tab name="gloo" codelang="shell" >}}
? What language do you wish to use for the filter:
    assemblyscript
  ▸ cpp
? With which platform do you wish to use the filter?:
    istio 1.5.x
  ▸ gloo 1.3.x
{{< /tab >}}
{{< tab name="istio" codelang="shell">}}
? What language do you wish to use for the filter:
    assemblyscript
  ▸ cpp
? With which platform do you wish to use the filter?:
  ▸ istio 1.5.x
    gloo 1.3.x
{{< /tab >}}
{{< /tabs >}}

```
INFO[0014] extracting 5072 bytes to /Users/ilackarms/go/src/github.com/solo-io/wasme/cpp-filter
```

The `init` command will place our *base* filter into the `cpp-filter` directory:

```shell
cd cpp-filter
tree .
```

```

├── BUILD
├── README.md
├── WORKSPACE
├── bazel
│   └── external
│       ├── BUILD
│       ├── emscripten-toolchain.BUILD
│       └── envoy-wasm-api.BUILD
├── filter.cc
├── filter.proto
├── runtime-config.json
└── toolchain
    ├── BUILD
    ├── cc_toolchain_config.bzl
    ├── common.sh
    ├── emar.sh
    └── emcc.sh
```

`wasme` uses [Bazel](https://bazel.build/) to build C++ filters under the hood.

{{% notice note %}}
The `runtime-config.json` file present in WASM filter modules is required by `wasme` to build the filter.

At least must one valid [`root_id`](https://github.com/envoyproxy/envoy-wasm/blob/master/api/envoy/config/wasm/v2/wasm.proto#L47)
matching the WASM Filter must be present in the `rootIds` field.
{{% /notice %}}

## Making changes to the base filter

The new directory contains all files necessary to build and deploy a WASM filter with `wasme`. A brief description of each file is found below:

| File | Description |
| ----- | ---- |
| `BUILD`                | The Bazel BUILD file used to build the filter. |         
| `WORKSPACE`            | The Bazel WORKSPACE file used to build the filter. |         
| `bazel/`               | Bazel external dependencies. |              
| `toolchain/`           | Bazel tooling for building wasm modules. |              
| `filter.cc`            | The source code for the filter, written in C++. |         
| `filter.proto`         | The protobuf schema of the filter configuration. |         
| `runtime-config.json`  | Config stored with the filter image used to load the filter at runtime. |

Open `filter.cc` in your favorite text editor. We'll make some changes to customize our new filter.

Navigate to the `AddHeaderContext::onResponseHeaders` method defined near the bottom of the file.
 Let's add a new header that we can use to verify our module was executed correctly. Let's add a new response header `hello: world!`:

```typescript
    addResponseHeader("hello", "world!");
```

Your method should look like this:

```c++
FilterHeadersStatus AddHeaderContext::onResponseHeaders(uint32_t) {
  addResponseHeader("hello", "world!");
  return FilterHeadersStatus::Continue;
}

```

The code above will add the `hello: world!` header to HTTP responses processed by our filter.

## Building the filter

Now, let's build a WASM image from our filter with `wasme`. The filter will be tagged and stored in a local registry, similar to how [Docker](https://www.docker.com/) stores images. 

Images tagged with `wasme` have the following format:

```
<registry address>/<registry username|org>/<image name>:<version tag>
```

* `<registry address>` specifies the address of the remote OCI registry where the image will be pushed by the `wasme push` command. The project authors maintain a free public registry at `webassemblyhub.io`.
 
* `<registry username|org>` either your username for the remote OCI registry, or a valid org name with which you are registered.


*See the [`wasme push`]({{< versioned_link_path fromRoot="/tutorial_code/push_tutorials">}}) documentation for instructions on pushing filters built with `wasme`.*


In this example we'll include the registry address `webassemblyhub.io` so our image can be pushed to the remote registry, along with GitHub username which will be used to authenticate to the registry.

Build and tag our image like so:

```shell
wasme build . -t webassemblyhub.io/<USERNAME>/add-header:v0.1
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
