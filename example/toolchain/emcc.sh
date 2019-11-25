#!/bin/bash
set -euo pipefail

. $(dirname $0)/common.sh

emcc -s EMIT_EMSCRIPTEN_METADATA=1 -s STANDALONE_WASM=1 -s EXPORTED_FUNCTIONS=['_malloc','_free'] "$@"

# Remove the first line of .d file
# not sure why... https://docs.bazel.build/versions/master/tutorial/cc-toolchain-config.html
# also, sorten the prefix as it seems that our clang doesn't support -fno-canonical-system-headers
find . -name "*.d" -exec sed -i '' -e '2d' -e 's%[^ ]*/external/emscripten_toolchain/upstream/emscripten/system/%external/emscripten_toolchain/upstream/emscripten/system/%' {} \;

# yet another hack till i can figure out how to make no-canonical-prefixes work
find . -name "*.d" -exec sed -i '' -e 's%[^ ]*/external/envoy_wasm_api-tmp/api/wasm/cpp%external/envoy_wasm_api%' {} \;
