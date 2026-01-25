#!/bin/bash

## Runs all the TypeScript tests. Pass "-c" to clear the cache first, "-w" to watch for changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

do_clean=false
do_watch=false

while getopts cw option
do
  case "${option}" in
    c) do_clean=true;;
    w) do_watch=true;;
    *) echo "unknown option"; exit 1;;
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
