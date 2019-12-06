#!/usr/bin/env sh

# must mount /tmp/build_output
DESTDIR=/tmp/build_output

bazel build :filter.wasm

cp -r bazel-bin/filter.wasm $DESTDIR/
