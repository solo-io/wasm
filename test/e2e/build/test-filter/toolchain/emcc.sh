#!/bin/bash
set -euo pipefail

. $(dirname $0)/common.sh

emcc -s EMIT_EMSCRIPTEN_METADATA=1 -s STANDALONE_WASM=1 -s EXPORTED_FUNCTIONS=['_malloc','_free'] "$@"

# clang doesn't support `-no-canonical-system-headers` so sed it
# find the .d file in the args and fix it:

for arg in "$@"
do
    if [ ${arg: -2} == ".d" ]; then
        echo Fixing $arg
        sed -e 's%[^ ]*/external/emscripten_toolchain/upstream/emscripten/system/%external/emscripten_toolchain/upstream/emscripten/system/%' $arg > $arg.tmp
        mv $arg.tmp $arg
        # some zlib headers are treated as system headers
        sed -e 's%[^ ]*/external/zlib/%external/zlib/%' $arg > $arg.tmp
        mv $arg.tmp $arg
        break
    fi
done
