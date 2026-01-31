#!/bin/bash

## Full release build: client deps, assets, templates, and binary.
##
## Usage:
##   ./bin/build/fullbuild.sh
##
## Requires:
##   - Node.js and npm
##   - Go toolchain
##   - make
##
## Notes:
##   - Runs npm install, builds client assets, compiles templates,
##     runs make clean, go mod tidy, then builds the release binary.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd npm "install Node.js from https://nodejs.org"
require_cmd go "install Go from https://go.dev/dl/"
require_cmd make "install Xcode Command Line Tools or build-essential"

cd client
npm install
cd ..
./bin/build/client.sh

make clean
./bin/templates.sh
go mod tidy
./bin/build/build.sh
