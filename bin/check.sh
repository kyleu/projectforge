#!/bin/bash

## Runs code statistics, checks for outdated dependencies, then runs linters

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

echo "=== linting ==="
golangci-lint -c ./tools/release/.golangci.yml run --fix --max-issues-per-linter=0 --sort-results ./...
