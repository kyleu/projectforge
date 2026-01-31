#!/bin/bash

## Generates a Go module dependency graph SVG.
##
## Usage:
##   ./bin/util/view-go-deps.sh
##
## Requires:
##   - gomod in PATH
##   - graphviz in PATH
##   - dot in PATH
##
## Outputs:
##   - ./tmp/deps.svg

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd gomod "install a gomod CLI that supports 'gomod graph'"
require_cmd dot "install graphviz (dot)"

echo "building dependency SVG..."
gomod graph | dot -Tsvg -o ./tmp/deps.svg
echo "Output written to ./tmp/deps.svg"

open ./tmp/deps.svg
