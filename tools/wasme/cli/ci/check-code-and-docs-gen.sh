#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
  echo "_output" > tools/wasme/cli/.gitignore
fi

make install-deps

set +e

make manifest-gen generated-code -B > /dev/null
if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi

# Tars can build slightly differently based on the host system.
# At build time, they will be re-built in CI, so we can ignore diffs in tars
git status --porcelain | grep archive_2gobytes.go | cut -c 20- | xargs git restore

if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Error: Generating code produced a non-empty diff"
  echo "Try running 'make clean install-deps manifest-gen generated-code -B' from the tools/wasme/cli directory, then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi
