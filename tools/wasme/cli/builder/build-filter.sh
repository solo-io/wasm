#!/usr/bin/env sh

set -e

ls /build_output || echo "/build_output must be mounted to this container" exit 1

# must have mounted /build_output
DESTFILE=${DESTFILE:-/build_output/filter.wasm}

bazel_build() {
  # subpath within bazel-bin where bazel puts the file
  BAZEL_OUTPUT=${BAZEL_OUTPUT:-filter.wasm}

  # name of bazel target to run
  TARGET=${TARGET:-:filter.wasm}

  echo running "bazelisk build ${TARGET}"
  bazelisk build ${TARGET}

  cp -r bazel-bin/${BAZEL_OUTPUT} ${DESTFILE}
}

npm_build() {
  # subpath within source dir where node puts the file
  NPM_OUTPUT=${NPM_OUTPUT:-build/optimized.wasm}

  if [ -z $NPM_USERNAME ]; then
      echo skipping login
    else
      echo running "creating npm user $NPM_USERNAME"
      /usr/bin/expect <<EOD
spawn npm adduser
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD
  fi

  echo running "npm install && npm run asbuild"
  npm install && npm run asbuild

  cp -r ${NPM_OUTPUT} ${DESTFILE}

  rm -rf build
}

tinygo_build() {
  tinygo build -o ${DESTFILE} -target=wasi -wasm-abi=generic .
}

echo -n "Building with ${BUILD_TOOL}..."

case $BUILD_TOOL in

  npm)
    npm_build
    ;;

  bazel)
    bazel_build
    ;;
  tinygo)
    tinygo_build
    ;;
  *)
    echo -n "unsupported build tool: ${BUILD_TOOL}"
    exit 1
    ;;
esac
