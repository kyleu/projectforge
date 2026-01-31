#!/bin/bash

## Builds the WASM server assets and binary.
##
## Usage:
##   ./bin/build/wasmserver.sh
##
## Requires:
##   - Go toolchain
##   - Built client assets in ./assets (for maps/icons)
##
## Outputs:
##   - build/wasm/{{{ .Exec }}}.wasm
##   - build/wasm/* (html/js/sw/assets)

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

mkdir -p ./build/wasm/assets

cp ./tools/wasmserver/index.html ./build/wasm/index.html
cp ./tools/wasmserver/server.js ./build/wasm/server.js
cp ./tools/wasmserver/sw.js ./build/wasm/sw.js
cp ./tools/wasmserver/wasm_exec.js ./build/wasm/wasm_exec.js

cp ./assets/client.css.map ./build/wasm/assets/client.css.map
cp ./assets/client.js.map ./build/wasm/assets/client.js.map
cp ./assets/favicon.ico ./build/wasm/favicon.ico
cp ./assets/logo.svg ./build/wasm/assets/logo.svg
cp ./assets/logo.svg ./build/wasm/logo.svg

echo "building {{{ .Name }}} WASM server library..."
GOOS=js GOARCH=wasm go build -o ./build/wasm/{{{ .Exec }}}.wasm .
echo "Output written to ./build/wasm ({{{ .Exec }}}.wasm and assets)"
