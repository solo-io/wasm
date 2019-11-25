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

source external/emscripten_toolchain/emsdk_env.sh

 # the emscripten sdk does some path comparison, so make EM_CACHE an absolute path to make it work. 
mkdir -p "tmp/emscripten_cache"
export EM_CACHE=${PWD}"/tmp/emscripten_cache"
export TEMP_DIR="tmp"