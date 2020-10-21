#!/bin/bash

set -ex

protoc --version

if [ ! -f .gitignore ]; then
  echo "_output" > .gitignore
fi

make install-deps

set +e

make generated-code operator-gen manifest-gen -B > /dev/null
if [[ $? -ne 0 ]]; then
  echo "Code generation failed"
  exit 1;
fi
if [[ $(git status --porcelain | wc -l) -ne 0 ]]; then
  echo "Error: Generating code produced a non-empty diff"
  echo "Try running 'make install-go-tools generated-code operator-gen manifest-gen -B' from the tools/wasme/cli directory, then re-pushing."
  git status --porcelain
  git diff | cat
  exit 1;
fi
