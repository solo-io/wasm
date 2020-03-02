# Compiling rust to .wasm files using the SDK

The rust SDK is a WIP.

## rust if you want to use it or rebuild the rust WebAssembly tests

`curl https://sh.rustup.rs -sSf | sh`

## Add the WASM target
Need to add the wasm target before we can compile rust for it:

`rustup target add wasm32-unknown-unknown`

## Building WASM
The example projects use cargo, so there are two ways to build the project:

1. Will build a development copy:
    `cargo build --target wasm32-unknown-unknown`
2. Build for release:
    `cargo build --target wasm32-unknown-unknown --release`

## Deploy using wasme
Now that the build is ready you need to package it into an OCI image. First, navigate to
the `rust/target/wasm32-unknown-unknown/release` directory.

You need to tell `wasme` some details about your filter. Change the `config.json` file
to reflect your filters details. Most importantly you want to change the filter `rootId` field.

Now you need to actually generate the image. You can do so with `wasme` by running:
`wasme build precompiled target/wasm32-unknown-unknown/release/wasm_filter_bindings.wasm --tag webassemblyhub.io/<github username>/<filter-name>:v0.5 --config config-13.json`

Once the image is built you can publish it to the Hub using `wasme`. Run:
`wasme push webassemblyhub.io/<github username>/<filter-name>:v0.1 ./wasm_filter_bindings.wasm`

## More about wasme
For more information about using wasme and the WebAssemblyHub [check here](https://docs.solo.io/web-assembly-hub/latest/tutorial_code/getting_started_1/).

