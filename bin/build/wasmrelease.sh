#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.
## Builds the application as a WebAssembly library

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="0.0.0"

./bin/build/wasmserver.sh

cd "build/wasm"
zip -r "../dist/projectforge_${TGT}_wasm_html.zip" ./*
