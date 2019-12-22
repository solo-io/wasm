#!/usr/bin/env sh

set -ex

# must have mounted /tmp/build_output
DESTFILE=${DESTFILE:-/tmp/build_output/filter.wasm}

# subpath within bazel-bin where bazel puts the file
BAZEL_OUTPUT=${BAZEL_OUTPUT:-filter.wasm}

# name of bazel target to run
TARGET=${TARGET:-:filter.wasm}

echo running "bazel build ${TARGET}"
bazel build ${TARGET}

cp -r bazel-bin/${BAZEL_OUTPUT} ${DESTFILE}

bazel clean
