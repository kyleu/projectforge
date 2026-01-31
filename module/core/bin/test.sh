#!/bin/bash

## Runs Go tests with gotestsum.
##
## Usage:
##   ./bin/test.sh [-c|--clean] [-w|--watch]
##
## Arguments:
##   -c, --clean  Clear the Go test cache before running.
##   -w, --watch  Watch mode (passes --watch to gotestsum).
##
## Environment:
##   - Loads env vars from ./test.env if present.
##
## Requires:
##   - Go toolchain
##   - gotestsum in PATH
##
## Notes:
##   - Runs ./bin/test-setup.sh if it exists.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "Go toolchain"
require_cmd gotestsum "install via 'go install gotest.tools/gotestsum@latest'"

do_clean=false
do_watch=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    -c|--clean) do_clean=true; shift;;
    -w|--watch) do_watch=true; shift;;
    --) shift; break;;
    *) echo "unknown option: $1" >&2; exit 1;;
  esac
done

test_args=""

if $do_clean; then
  echo "cleaning test cache...";
  go clean -testcache
fi

if $do_watch; then
  echo "watching for file changes...";
  test_args="${test_args} --watch"
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

gotestsum${test_args} -- -race ./app/...
