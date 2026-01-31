#!/bin/bash

## Builds client assets using esbuild/tsc via build.js.
##
## Usage:
##   ./bin/build/client.sh
##
## Requires:
##   - Node.js and npm
##   - esbuild and tsc available to build.js

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../../client"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd node "install Node.js from https://nodejs.org"
require_cmd npm "install Node.js from https://nodejs.org"

node build.js
echo "Output written to [../assets]"
