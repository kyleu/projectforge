#!/bin/bash

## Runs Go linters via golangci-lint.
##
## Usage:
##   ./bin/check.sh
##
## Requires:
##   - golangci-lint in PATH

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd golangci-lint "install from https://golangci-lint.run/usage/install/"

echo "=== linting ==="
golangci-lint run --max-issues-per-linter=0 ./...
