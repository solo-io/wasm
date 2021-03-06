# Wasm Image specifications

## Introduction

The Wasm Image specification defines how to bundle Wasm modules as container images. A compatible Wasm image consists of a Wasm binary file, and runtime metadata for the target Wasm runtime. We primarily consider the use case of [Envoy] as Wasm runtime and its Wasm filter/plugins, although the spec is intended to be generic and to provide a standard mechanism to manage the building and running of Wasm modules on any Wasm runtime.

## Terminology:

| Term                               | Definition                                       |
|------------------------------------|--------------------------------------------------|
| Wasm Module                        | The distributable, loadable, and executable unit of code in WebAssembly. 
| Wasm Image Specification           | The specification for storing Wasm modules as container images.
| Wasm Runtime                       | The execution environment into which a Wasm Module may be loaded. This refers to the application itself which loads and executes a wasm module. Examples include web browsers, the Open Policy Agent, the Envoy Proxy, or any other application which supports extension via Wasm modules. 
| Runtime Configuation              | Configuration or metadata specific to the runtime which consumes a module. 

## Specifications

Here we have several specifications for how to bundle Wasm modules as container images. 

There are two variants of the specification:
- [spec.md](spec.md)
- [spec-compat.md](spec-compat.md)

Developers and Wasm module consumers can leverage both of these specifications. 

For clarity, we call the variant in [spec.md](spec.md) *oci*, and the one in [spec-compat.md](spec-compat.md) *compat*.

## Difference between variants

Our goal is to make the *oci* variant the default format for shipping Wasm modules in container images, we acknowledge however that there are toolchains and registries deployed and in use that do not support our custom media types yet. To accomodate those exisiting toolchains, there is the semantically equivalent *compat* variant, which provides the same feature set, but is compatible with existing tooling because it 'looks' very much like a normal container image. Implementations of this spec should support both variants.

With that said, the key difference between these two variants is that, the *oci* variant makes use of the custom media types on [OCI Artifact] for image layers while the *compat* variant leverages the standard media types.

As a consequence, that introduces the difference in tools available for building and pushing images. 
For example, the only way to build and push *oci* variant images is to use [`wasme`] cli while you can use [`buildah`] or [`docker`] for *compat* variant images.

Not only that, the usage of custom media types on top of [OCI Artifact] requires registries to support arbitrary custom media types. Therefore, you might not be able to push *oci* variants to your registry while [WebAssemblyHub] is designed for accepting them.

## Wasm image support in [Istio]

Istio's Wasm Plugin API has support for **both of these variants** to deploy your Wasm plugins into Envoy sidecars.

### Which variant should I use in [Istio]?

Given that Istio supports both of variants, you can choose whichever variant depending on your needs. For example, if you want to use [`docker`] cli, then choose *compat* variant and push them to your container registries. You might want to use [`wasme`] cli and [WebAssemblyHub] then choose the *oci* variant.

## How can I build images?

For *oci* variant, see the guideline by [`wasme`].

For *compat* variant, follow the instructions in 
- [build a compat image with buildah](spec-compat.md#appendix-1-build-a-compat-image-with-buildah) if you want to use [`buildah`].
- [build a compat image with docker](spec-compat.md#appendix-2-build-a-compat-image-with-docker-cli) if you want to use [`docker`].


[Envoy]: https://github.com/envoyproxy/envoy
[Istio]: https://github.com/istio/istio
[OCI Artifact]: https://github.com/opencontainers/artifacts
[WebAssemblyHub]: https://webassemblyhub.io/

[`docker`]: https://docs.docker.com/engine/reference/commandline/cli/
[`buildah`]: https://github.com/containers/buildah
[`wasme`]: https://docs.solo.io/web-assembly-hub/latest/installation/
