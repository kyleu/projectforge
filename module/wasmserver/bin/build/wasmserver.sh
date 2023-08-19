#!/bin/bash

## Builds the application as a WebAssembly library

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

mkdir -p ./build/wasm
cp ./tools/wasmserver/index.html ./build/wasm/index.html
cp ./tools/wasmserver/wasm_exec.js ./build/wasm/wasm_exec.js
cp ./assets/client.js ./build/wasm/client.js
cp ./assets/client.css ./build/wasm/client.css
cp ./assets/logo.svg ./build/wasm/logo.svg

echo "building Project Forge WASM server library..."
GOOS=js GOARCH=wasm go build -o ./build/wasm/{{{ .Exec }}}.wasm .
