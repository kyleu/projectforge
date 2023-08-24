#!/bin/bash

## Starts the app in wasm mode, please run `./bin/build/wasmserver.sh` before starting this.

## Alternately, run `air -c tools/wasmserver/.air.conf` from the project root to reload on changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

npx http-server build/wasm --port 40009 --ext html --gzip --brotli --cors
