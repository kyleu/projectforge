#!/bin/bash

## Builds a darwin/arm64 release binary and runs a single action.
##
## Usage:
##   ./bin/export.sh [action]
##
## Arguments:
##   action  Project Forge action to run (default: preview).
##
## Requires:
##   - Go toolchain
##   - make
##
## Notes:
##   - Uses ./bin/build/build.sh darwin arm64.
##   - Executes ./build/darwin/arm64/projectforge <action>.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd make "install Xcode Command Line Tools or build-essential"

act=${1:-preview}
"${dir}"/../bin/build/build.sh darwin arm64 && "${dir}"/../build/darwin/arm64/projectforge "$act"
