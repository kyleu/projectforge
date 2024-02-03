#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the application as a WebAssembly library

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

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

echo "building Project Forge WASM server library..."
GOOS=js GOARCH=wasm go build -o ./build/wasm/projectforge.wasm .
