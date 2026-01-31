#!/bin/bash

## Builds the release binary and writes symbol sizes to ./tmp/nm.txt.
##
## Usage:
##   ./bin/util/view-binary-size.sh
##
## Requires:
##   - Go toolchain
##   - c++filt (for demangling)
##   - make
##
## Outputs:
##   - ./tmp/nm.txt

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
require_cmd c++filt "install binutils for c++filt"

make build-release
go tool nm -size build/release/{{{ .Exec }}} | c++filt > ./tmp/nm.txt
echo "Output written to ./tmp/nm.txt"
