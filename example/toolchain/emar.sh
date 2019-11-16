#!/bin/bash
# emsdk_env.sh\emcc doesn't like the bazel sand box
if [[ "$OSTYPE" == "linux-gnu" ]]; then
cd -P /proc/self/cwd
fi

set -euo pipefail
source external/emscripten_toolchain/emsdk_env.sh

 # mkdir -p "tmp/emscripten_cache"
 # mkdir -p "tmp/tmp"
 # export EM_CACHE="tmp/emscripten_cache"
 # export TEMP_DIR="tmp/tmp"

export EM_CACHE="/tmp/emscripten_cache"

emar "$@"

# as this stands this build is not hermetic and will interfere with the environment installed outside
# fix TBD