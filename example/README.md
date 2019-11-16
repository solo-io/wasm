This compiles the example filter from envoy.
Currently the build system is not hermetic, but it is pretty easy to use.

build with
```
bazel build :envoy_filter_http_wasm_example.wasm --config=wasm --sandbox_writable_path $(bazel info output_base)/external/emscripten_toolchain/.emscripten_cache/
```