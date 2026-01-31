#!/bin/bash

## Serves the built WASM assets for local dev.
##
## Usage:
##   ./bin/wasmdev.sh
##
## Requires:
##   - Node.js and npx
##   - http-server (npx http-server)
##
## Notes:
##   - Run ./bin/build/wasmserver.sh first to generate ./build/wasm.
##   - Serves on http://localhost:40009 with gzip/brotli enabled.
##   - Alternately, run `air -c tools/wasmserver/.air.conf` from the project root to reload on changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd npx "install Node.js from https://nodejs.org"

npx http-server build/wasm --port 40009 --ext html --gzip --brotli --cors
