#!/bin/bash

# emsdk_env.sh\emcc doesn't like the bazel sandbox
# specifically, emsdk_env.sh seems to try to `cd` and `cd` back which doesn't work well
if [[ "$OSTYPE" == "linux-gnu" ]]; then
cd -P /proc/self/cwd
fi

export NODE_JS=''
export EMSCRIPTEN_ROOT='external/emscripten_toolchain'
export SPIDERMONKEY_ENGINE=''
export EM_EXCLUSIVE_CACHE_ACCESS=1
export EMCC_SKIP_SANITY_CHECK=1
export EMCC_WASM_BACKEND=1

export TEMP_DIR="tmp"

source external/emscripten_toolchain/emsdk_env.sh
