#!/bin/sh

set -eu

WASME_VERSIONS=$(curl -sH"Accept: application/vnd.github.v3+json" https://api.github.com/repos/solo-io/wasme/releases | python -c "import sys; from json import loads as l; releases = l(sys.stdin.read()); print('\n'.join(release['tag_name'] for release in releases))")

if [ "$(uname -s)" = "Darwin" ]; then
  OS=darwin
else
  OS=linux
fi

for WASME_VERSION in $WASME_VERSIONS; do

tmp=$(mktemp -d /tmp/wasme.XXXXXX)
filename="ap-${OS}-amd64"
url="https://github.com/solo-io/wasme/releases/download/${WASME_VERSION}/${filename}"

if curl -f ${url} >/dev/null 2>&1; then
  echo "Attempting to download Wasme CLI version ${WASME_VERSION}"
else
  continue
fi

(
  cd "$tmp"

  echo "Downloading ${filename}..."

  SHA=$(curl -sL "${url}.sha256" | cut -d' ' -f1)
  curl -sLO "${url}"
  echo "Download complete!, validating checksum..."
  checksum=$(openssl dgst -sha256 "${filename}" | awk '{ print $2 }')
  if [ "$checksum" != "$SHA" ]; then
    echo "Checksum validation failed." >&2
    exit 1
  fi
  echo "Checksum valid."
)

(
  cd "$HOME"
  mkdir -p ".wasme/bin"
  mv "${tmp}/${filename}" ".wasme/bin/wasme"
  chmod +x ".wasme/bin/wasme"
)

rm -r "$tmp"

echo "Wasme CLI was successfully installed 🎉"
echo ""
echo "Add the Wasme CLI to your path with:"
echo "  export PATH=\$HOME/.wasme/bin:\$PATH"
echo ""
echo "Now run:"
echo "  ap init myproject     # generate a new project directory"
echo "Please see visit the WebAssembly Hub guides for more:  https://docs.solo.io/web-assembly-hub/latest"
exit 0
done

echo "No versions of wasme found."
exit 1
