---
title: "Introduction"
weight: 1
---

<h1 align="center">
    <img src="https://github.com/solo-io/wasme/blob/master/docs/content/img/logo.png?raw=true" alt="WebAssembly Hub" width="371" height="242">
  <br>
  The WebAssembly Hub
</h1>

The WebAssembly Hub is a meeting place for the community to share and consume WebAssembly Envoy extensions. Easily search and find extensions that meet the functionality you want to add and give them a try.

Please see the [announcement blog](https://medium.com/solo-io/introducing-the-webassembly-hub-a-service-for-building-deploying-sharing-and-discovering-wasm-d461719383ca) that goes into more detail on the motivation for WebAssembly Hub and how we see it driving the future direction of Envoy-based networking projects/products including API Gateways and Service Mesh.

The `wasme` CLI provides a tool for building and sharing Envoy WebAssembly extensions.

[**Installation**](https://docs.solo.io/web-assembly-hub/latest/installation/) &nbsp; |
&nbsp; [**Documentation**](https://docs.solo.io/web-assembly-hub/latest) &nbsp; |
&nbsp; [**Blog**](https://medium.com/solo-io/introducing-the-webassembly-hub-a-service-for-building-deploying-sharing-and-discovering-wasm-d461719383ca) &nbsp; |
&nbsp; [**Slack**](https://slack.solo.io) &nbsp; |
&nbsp; [**Twitter**](https://twitter.com/soloio_inc)

### How does it work?

The WebAssembly Hub, in combination with the `wasme` CLI, provides an easy way to build, push, pull, and share Envoy WebAssembly Filters.

The WebAssembly Hub acts as an image registry for WebAssembly Filters hosted at https://webassemblyhub.io. Use the `wasme` CLI to:

- compile [Envoy WebAssembly](https://github.com/envoyproxy/envoy-wasm) filters on a local machine (the only dependency is `docker`)
- push filters to https://webassemblyhub.io
- pull filters from https://webassemblyhub.io
- publish filters to the catalog at https://webassemblyhub.io/extensions/ 

### Getting Started

See the [Getting Started tutorial]({{< versioned_link_path fromRoot="/tutorial_code/getting_started">}}) to build, push, and run your first WebAssembly Filter!

### Next Steps
- Join us on our Slack channel: [https://slack.solo.io/](https://slack.solo.io/)
- Follow us on Twitter: [https://twitter.com/soloio_inc](https://twitter.com/soloio_inc)
- Check out the docs: [https://docs.solo.io/web-assembly-hub/latest](https://docs.solo.io/web-assembly-hub/latest)
- Contribute to the [Docs](https://github.com/solo-io/wasme)
