#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Starts the app in wasm mode, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

go mod tidy
ulimit -S -n 65536

./bin/build/wasmserver.sh
http -i --cors ./build/wasm
