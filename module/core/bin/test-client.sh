#!/bin/bash

## Runs client tests via npm.
##
## Usage:
##   ./bin/test-client.sh [-c|--clean] [-w|--watch]
##
## Arguments:
##   -c, --clean  Clear Vitest cache before running.
##   -w, --watch  Watch mode (passes --watch to the test runner).
##
## Environment:
##   - Loads env vars from ./test.env if present.
##
## Requires:
##   - Node.js and npm
##   - Client deps installed (npm install in ./client)
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

require_cmd npm "install Node.js from https://nodejs.org"

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

if $do_clean; then
  echo "cleaning test cache...";
  rm -rf "$dir/../client/node_modules/.vitest" "$dir/../client/node_modules/.cache/vitest"
fi

test_args=()

if $do_watch; then
  echo "watching for file changes...";
  test_args+=(--watch)
fi

cd "$dir/../client"

echo "=== testing client ==="
if [ ${#test_args[@]} -gt 0 ]; then
  npm run test -- "${test_args[@]}"
else
  npm run test
fi
