#!/bin/bash

## Builds and runs the desktop app for local development.
##
## Usage:
##   ./bin/build/desktop.sh
##
## Requires:
##   - Go toolchain and make
##
## Notes:
##   - Uses ./tools/desktop to build and runs build/debug/{{{ .Exec}}}-desktop.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd make "install Xcode Command Line Tools or build-essential"

cd tools/desktop

go mod tidy

make build
echo "build complete, starting desktop application"

cd "$dir/../.."
echo "Builds written to ./build/debug/{{{ .Exec }}}-desktop"
build/debug/{{{ .Exec }}}-desktop
