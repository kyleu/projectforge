#!/bin/bash

## Formats Go sources with gofumpt.
##
## Usage:
##   ./bin/format.sh
##
## Requires:
##   - gofumpt in PATH
##
## Notes:
##   - Excludes ./data, ./module, ./testproject, ./assets/module.
##   - Skips generated *.html.go and *.sql.go files.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd gofumpt "install via 'go install mvdan.cc/gofumpt@latest'"

echo "=== formatting ==="
gofumpt -w $(find . -type f -name "*.go" | grep -v \\./data | grep -v \\./module | grep -v \\./testproject | grep -v \\./assets/module | grep -v .html.go | grep -v .sql.go)
