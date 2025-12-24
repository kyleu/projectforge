#!/bin/bash

## Runs code statistics, checks for outdated dependencies, then runs linters

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

echo "=== linting ==="
golangci-lint-v2 run --max-issues-per-linter=0 ./...
