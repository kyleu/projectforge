#!/bin/bash

## Runs Go tests with race + coverage output.
##
## Usage:
##   ./bin/coverage.sh [-c|--clean]
##
## Arguments:
##   -c, --clean  Clear the Go test cache before running.
##
## Environment:
##   - Loads env vars from ./test.env if present.
##
## Requires:
##   - Go toolchain
##   - gotestsum
##   - go-cover-treemap (optional, for SVG output)
##
## Outputs:
##   - ./tmp/coverage.out
##   - ./tmp/coverage.svg (if go-cover-treemap is installed)
##
## Notes:
##   - Runs ./bin/test-setup.sh if it exists.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."
mkdir -p ./tmp

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "Go toolchain"
require_cmd gotestsum "install via 'go install gotest.tools/gotestsum@latest'"

do_clean=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    -c|--clean) do_clean=true; shift;;
    --) shift; break;;
    *) echo "unknown option: $1" >&2; exit 1;;
  esac
done

test_args=""

if $do_clean; then
  echo "cleaning test cache...";
  go clean -testcache
fi

if [ -f "test.env" ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    if [[ -n "$line" && ! $line =~ ^#.* ]]; then
      export "${line?}"
    fi
  done < "test.env"
fi

if [ -f "./bin/test-setup.sh" ]; then
	./bin/test-setup.sh
fi

gotestsum${test_args} -- -race -coverprofile ./tmp/coverage.out ./app/...

if command -v go-cover-treemap &> /dev/null; then
	echo "generating code coverage SVG..."
	go-cover-treemap -coverprofile ./tmp/coverage.out > ./tmp/coverage.svg
fi
