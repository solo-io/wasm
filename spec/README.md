# Wasm Image specifications

## Overview

Here we have several specifications for how to bundle Wasm modules as container images. 
They primarily consider the use case of [Envoy] Wasm plugins.

There are two variants of the specification:
- [spec.md](spec.md)
- [spec-compat.md](spec-compat.md)

Developers and Wasm module consumers can leverage both of these specifications. 

For clarity, we call the variant in [spec.md](spec.md) *oci*, and the one in [spec-compat.md](spec-compat.md) *compat*.

## Difference between variants

The key difference between these two variants is that, the *oci* variant makes use of the custom media types on [OCI Artifact] for image layers while the *compat* variant leverages the standard media types.

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
