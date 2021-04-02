---
title: "Building WASM Filters in AssemblyScript"
weight: 1
description: "Build a simple WebAssembly filter in AssemblyScript."
---

In this tutorial we will write an Envoy filter in [AssemblyScript](https://docs.assemblyscript.org/) and build it using `wasme`.

## Creating an AssemblyScript WASM module

Refer to the [installation guide]({{< versioned_link_path fromRoot="/installation">}}) for installing `wasme`, the WebAssembly Hub CLI.

Let's create a new filter called `assemblyscript-filter`:

```shell
wasme init assemblyscript-filter
```

You'll be asked with an interactive prompt which language platform you are building for. At time of writing, the AssemblyScript Filter base is compatible with gloo:1.3.x, gloo:1.5.x, gloo:1.6.x, istio:1.5.x, istio:1.6.x, istio:1.7.x, istio:1.8.x, istio:1.9.x:

You should get output like this:

```shell script
Use the arrow keys to navigate: ↓ ↑ → ← 
? What language do you wish to use for the filter: 
  ▸ cpp
    rust
    assemblyscript
    tinygo

✔ assemblyscript
Use the arrow keys to navigate: ↓ ↑ → ← 
? With which platforms do you wish to use the filter?: 
  ▸ gloo:1.3.x, gloo:1.5.x, gloo:1.6.x, istio:1.5.x, istio:1.6.x, istio:1.7.x, istio:1.8.x, istio:1.9.x

✔ assemblyscript
✔ gloo:1.3.x, gloo:1.5.x, gloo:1.6.x, istio:1.5.x, istio:1.6.x, istio:1.7.x, istio:1.8.x, istio:1.9.x
```

```
INFO[0118] extracting 1812 bytes to /Users/ilackarms/go/src/github.com/solo-io/wasm/new-filter
```

The `init` command will place our *base* filter into the `assemblyscript-filter` directory:

```shell
cd assemblyscript-filter
tree .
```

You should get output like this:

```
.
├── assembly
│   ├── index.ts
│   └── tsconfig.json
├── package.json
├── package-lock.json
└── runtime-config.json
```

{{% notice note %}}
The `runtime-config.json` file present in WASM filter modules is required by `wasme` to build the filter.

At least must one valid [`root_id`](https://github.com/envoyproxy/envoy-wasm/blob/master/api/envoy/config/wasm/v2/wasm.proto#L47)
matching the WASM Filter must be present in the `rootIds` field.
{{% /notice %}}

Open this project in your favorite IDE. The source code is [AssemblyScript](https://github.com/AssemblyScript/assemblyscript) (a subset of [Typescript](https://www.typescriptlang.org/)) and we'll make some changes to customize our new filter.


## Making changes to the base filter

The new directory contains all files necessary to build and deploy a WASM filter with `wasme`. A brief description of each file is found below:

| File | Description |
| ----- | ---- |
| `assembly/index.ts`        | The source code for the filter, written in AssemblyScript. |
| `assembly/tsconfig.json`   | Typescript config file (AssemblyScript is a subset of Typescript). |
| `package.json`             | Used by to import npm modules during build time. |
| `package-lock.json`        | Locked npm modules.  |
| `runtime-config.json`      | Config stored with the filter image used to load the filter at runtime. |

Open `assembly/index.ts` in your favorite text editor. The source code is AssemblyScript and we'll make some changes to customize our new filter.

Navigate to the `onResponseHeaders` method defined near the top of the file. Let's add a new header that we can use to verify our module was executed correctly (later down in the tutorial). Let's add a new response header `hello: world!`:

```typescript
      stream_context.headers.response.add("hello", "world!");
```

Your method should look like this:

```typescript
    onResponseHeaders(a: u32): FilterHeadersStatusValues {
        const root_context = this.root_context;
        if (root_context.configuration == "") {
          stream_context.headers.response.add("hello", "world!");
        } else {
          stream_context.headers.response.add("hello", root_context.configuration);
        }
        return FilterHeadersStatusValues.Continue;
      }
```

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
wasme build assemblyscript -t webassemblyhub.io/$YOUR_USERNAME/add-header:v0.1 .
```

{{% notice note %}}
`wasme build` runs a build container inside of Docker which may run into issues due to SELinux (on Linux environments). To disable, run `sudo setenforce 0`
{{% /notice %}}

The module will take up to a few minutes to build. In the background, `wasme` has launched a Docker container to run the necessary
build steps.

When the build has finished, you'll be able to see the image with `wasme list`:

```bash
wasme list
```

You should get output like this:

```
NAME                                      TAG  SIZE    SHA      UPDATED
webassemblyhub.io/ilackarms/add-header v0.1 12.6 kB 0295d929 02 Apr 21 13:06 CST
```

## Next Steps

Now that we've successfully built our image, we can try [running it locally]({{< versioned_link_path fromRoot="/tutorial_code/deploy_tutorials/deploying_with_local_envoy">}}) or [pushing it to a remote registry]({{< versioned_link_path fromRoot="/tutorial_code/push_tutorials">}}) so it can be pulled and deployed in a Kubernetes environment.
