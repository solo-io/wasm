# Web Assembly

[Web Assembly][wasm] (WASM) is the future of cloud-native infrastructure extensibility. 

WASM is a safe, secure, and dynamic way of extending infrastructure with the language of your choice. [WASM tool chains][wasm-toolchain] compile your code from any of the [supported languages][supported-lang] into a type-safe, binary format that can be loaded dynamically in a WASM sandbox/VM. 

In this repo, you will find [tooling](./tools), [SDKs](./sdks), an [OCI-compatible](./spec) specification, and examples for working with WASM and specifically WASM on [Envoy Proxy][envoy] based frameworks (like [Gloo API Gateway][gloo] or [Istio Service Mesh][istio] -- but not limited to those)

One of those projects for working with WASM and Envoy proxy is [Web Assembly Hub][web-assembly-hub].

<h1 align="center">
    <img src="https://github.com/solo-io/wasme/blob/master/docs/content/img/logo.png?raw=true" alt="WebAssembly Hub" width="371" height="242">
</h1>

[WebAssembly Hub][web-assembly-hub] is a meeting place for the community to share and consume WebAssembly Envoy extensions. Easily search and find extensions that meet the functionality you want to add and give them a try.

Please see the [announcement blog][announcement] that goes into more detail on the motivation for WebAssembly Hub and how we see it driving the future direction of Envoy-based networking projects/products including API Gateways and Service Mesh.

# In this Repo

## Specification

In the [/spec](./spec) folder of this repo, you'll find an [OCI image][oci] specification for Web Assembly modules. This specification is an extension to OCI and describes how a WASM module is packaged with the appropriate metadata so that it can be distributed and loaded into Envoy-based frameworks/service meshes. For example, some of the metadata include the root-id of the module, the ABI versions that are targeted in Envoy, and basic name of the module. Please see the [spec](./spec) for more details.

## SDKs

In the [/sdks](./sdks) folder of this repo, you'll find a listing of the available SDKs for building Envoy Proxy based WASM modules. These SDKs implement the Envoy hooks for Envoy extension points in a language-specific way. The following exist today (with more on the way)

* C++
* Rust
* AssemblyScript
* TinyGo

## Tools

In the [/tools](./tools) folder, you'll find a set of tools for working with WASM modules. Specifically, tools to do the following:

* bootstrap a new WASM project with boilerplate/version alignment/tool chain 
* test that a WASM module can correctly be loaded into Envoy
* build your project with the specific tool chain
* publish your WASM module as a WASM OCI image to an OCI-compatible registry or [WebAssembly Hub][web-assembly-hub]
* load your WASM module into an Envoy-based framework like [Gloo][gloo] or [Istio][istio]

# Other Resources

* [Extending Envoy Proxy with WebAssembly](https://www.solo.io/blog/webinar-recap-extending-envoy-with-web-assembly/)
* [WebAssembly Hub Announcement][announcement]
* [WebAssembly Hub and Istio][announcement-istio]
* [Envoy WASM repo](https://github.com/envoyproxy/envoy-wasm)
* [Envoy WASM ABI](https://github.com/proxy-wasm/spec)
* [Proxy WASM repo](https://github.com/proxy-wasm)

# Community

The WASM projects in this repository are made up of a collaboration between Solo.io and the wider Web Assembly community. 

Please join the [Solo.io Slack @ #web-assembly-hub](https://slack.solo.io/) to participate!

[wasm]: http://webassembly.org
[envoy]: http://envoyproxy.io
[wasm-toolchain]: https://developer.mozilla.org/en-US/docs/WebAssembly/C_to_wasm#Emscripten_Environment_Setup
[supported-lang]: https://github.com/appcypher/awesome-wasm-langs
[web-assembly-hub]: https://webassemblyhub.io
[gloo]: https://gloo.solo.io
[istio]: https://istio.io
[oci]: https://github.com/opencontainers/image-spec/blob/master/spec.md
[announcement]: https://www.solo.io/blog/introducing-the-webassembly-hub-a-service-for-building-deploying-sharing-and-discovering-wasm/
[announcement-istio]: https://www.solo.io/blog/an-extended-and-improved-webassembly-hub-to-helps-bring-the-power-of-webassembly-to-envoy-and-istio/