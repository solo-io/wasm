#!/usr/bin/env sh

# must mount /tmp/build_output
DESTFILE=${DESTFILE:-filter.wasm}
BUILD_BASE=${BUILD_BASE:-.}
TARGET=${BUILD_BASE}:filter.wasm

echo running "bazel build ${TARGET}"
bazel build ${TARGET}

cp -r bazel-bin/${BUILD_BASE}/filter.wasm /tmp/build_output/${DESTFILE}

bazel clean
