#!/bin/bash
# emsdk_env.sh\emcc doesn't like the bazel sand box
cd -P /proc/self/cwd

set -euo pipefail
source external/emscripten_toolchain/emsdk_env.sh

 # mkdir -p "tmp/emscripten_cache"
 # mkdir -p "tmp/tmp"
 # export EM_CACHE="tmp/emscripten_cache"
 # export TEMP_DIR="tmp/tmp"

export EM_CACHE="/tmp/emscripten_cache"

emcc -s EMIT_EMSCRIPTEN_METADATA=1 -s STANDALONE_WASM=1 -s EXPORTED_FUNCTIONS=['_malloc','_free'] "$@"

# remove .d files so bazel wont be annoying
find . -name "*.d" -exec truncate -s 0 '{}' \;

# as this stands this build is not hermetic and will interfere with the environment installed outside
# fix TBD