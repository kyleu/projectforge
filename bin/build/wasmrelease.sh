#!/bin/bash

## Builds the WASM assets and zips them for release.
##
## Usage:
##   ./bin/build/wasmrelease.sh [version]
##
## Arguments:
##   version  Version tag for output filename (default: 0.0.0).
##
## Requires:
##   - Go toolchain
##   - zip
##
## Outputs:
##   - build/dist/projectforge_<version>_wasm_html.zip
##
## Notes:
##   - Runs ./bin/build/wasmserver.sh to prepare ./build/wasm.

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd zip "install zip from your package manager"

TGT=$1
[ "$TGT" ] || TGT="0.0.0"

./bin/build/wasmserver.sh

cd "build/wasm"
zip -r "../dist/projectforge_${TGT}_wasm_html.zip" ./*

cd "$dir/../.."
echo "Output written to ./build/dist/projectforge_${TGT}_wasm_html.zip"
