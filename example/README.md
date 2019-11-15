This compiles the example filter from envoy.
Currently the build system is not hermetic, but it is pretty easy to use.

build with
```
bazel build :yuval --config=wasm
```