#!/bin/bash

## Formatting code from all projects

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

echo "=== formatting ==="
gofumpt -w $(find . -type f -name "*.go" | grep -v \\./data | grep -v \\./module | grep -v \\./testproject | grep -v .html.go | grep -v .sql.go)
