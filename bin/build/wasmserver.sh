#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the application as a WebAssembly library

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

mkdir -p ./build/wasm
cp ./tools/wasmserver/index.html ./build/wasm/index.html
cp ./tools/wasmserver/wasm_exec.js ./build/wasm/wasm_exec.js
cp ./tools/wasmserver/server.js ./build/wasm/server.js
cp ./tools/wasmserver/sw.js ./build/wasm/sw.js
cp ./assets/logo.svg ./build/wasm/logo.svg

echo "building Project Forge WASM server library..."
GOOS=js GOARCH=wasm go build -o ./build/wasm/projectforge.wasm ./tools/wasmserver
