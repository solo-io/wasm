#!/bin/bash
set -euo pipefail


# emsdk_env.sh\emcc doesn't like the bazel sandbox
# specifically, emsdk_env.sh seems to try to `cd` and `cd` back which doesn't work well
cd -P /proc/self/cwd

export NODE_JS=''
export EMSCRIPTEN_ROOT='external/emscripten_toolchain'
export SPIDERMONKEY_ENGINE=''
export EM_EXCLUSIVE_CACHE_ACCESS=1
export EMCC_SKIP_SANITY_CHECK=1
export EMCC_WASM_BACKEND=1

 mkdir -p "tmp/emscripten_cache"
# the emscripten sdk does some path comparison, so make EM_CACHE an absolute path to make it work.
export EM_CACHE=$PWD/"tmp/emscripten_cache"
export TEMP_DIR="tmp"


source external/emscripten_toolchain/emsdk_env.sh


emcc -s EMIT_EMSCRIPTEN_METADATA=1 -s STANDALONE_WASM=1 -s EXPORTED_FUNCTIONS=['_malloc','_free'] "$@"


# Remove the first line of .d file
# not sure why... https://docs.bazel.build/versions/master/tutorial/cc-toolchain-config.html
# also, sorten the prefix as it seems that our clang doesn't support -fno-canonical-system-headers
find . -name "*.d" -exec sed -i -e '2d' -e 's%\S*/external/emscripten_toolchain/upstream/emscripten/system/%external/emscripten_toolchain/upstream/emscripten/system/%' {} \;

# yet another hack till i can figure out how to make no-canonical-prefixes work
find . -name "*.d" -exec sed -i -e 's%\S*/external/envoy_wasm-tmp/api/wasm/cpp%external/envoy_wasm%' {} \;
