#!/usr/bin/env sh

set -ex

WORKSPACE=$1
OUTDIR=$2

mkdir -p $OUTDIR

docker run \
    -v "${WORKSPACE}:/src/workspace" \
    -v "${OUTDIR}:/tmp/build_output" \
    -w /src/workspace \
    soloio/wasm-builder