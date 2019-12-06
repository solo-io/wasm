#!/usr/bin/env sh

WORKSPACE=$1
OUTDIR=$2

mkdir -p $OUTDIR

mkdir -p ./build_output && \
	docker run \
      -v "${WORKSPACE}:/src/workspace" \
      -v "${OUTDIR}:/tmp/build_output" \
      -w /src/workspace \
      soloio/wasm-builder