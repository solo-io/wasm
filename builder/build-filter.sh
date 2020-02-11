#!/usr/bin/env sh

set -ex

find .
pwd

ls /build_output || exit 1
echo I have /build_output!

# must have mounted /build_output
DESTFILE=${DESTFILE:-/build_output/filter.wasm}

bazel_build() {
  # subpath within bazel-bin where bazel puts the file
  BAZEL_OUTPUT=${BAZEL_OUTPUT:-filter.wasm}

  # name of bazel target to run
  TARGET=${TARGET:-:filter.wasm}

  echo running "bazel build ${TARGET}"
  bazel build ${TARGET}

  cp -r bazel-bin/${BAZEL_OUTPUT} ${DESTFILE}
}

npm_build() {
  # subpath within source dir where node puts the file
  NPM_OUTPUT=${NPM_OUTPUT:-build/optimized.wasm}

  echo running "npm install && npm run asbuild"
  npm install && npm run asbuild

  cp -r ${NPM_OUTPUT} ${DESTFILE}
}

echo -n "Building with ${BUILD_TOOL}..."

case $BUILD_TOOL in

  npm)
    npm_build
    ;;

  bazel)
    bazel_build
    ;;

  *)
    echo -n "unsupported build tool: ${BUILD_TOOL}"
    exit 1
    ;;
esac
