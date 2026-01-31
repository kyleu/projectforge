#!/bin/bash

## Builds a debug binary for the requested target.
##
## Usage:
##   ./bin/build/build-debug.sh [os] [arch] [filename]
##
## Arguments:
##   os        Target GOOS (default: darwin).
##   arch      Target GOARCH (default: amd64).
##   filename  Output binary name (default: projectforge).
##
## Requires:
##   - Go toolchain
##   - make
##
## Outputs:
##   - ./build/<os>/<arch>/<filename>

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

os=${1:-darwin}
arch=${2:-amd64}
fn=${3:-projectforge}

echo "Building debug $fn [$os $arch]..."
env GOOS="$os" GOARCH="$arch" make build
mkdir -p "./build/$os/$arch"
mv "./build/debug/$fn" "./build/$os/$arch/$fn"
echo "Executable written to [./build/$os/$arch/$fn]"
