#!/bin/bash

## Runs all the tests. Pass "-c" to clear the cache first, "-w" to watch for changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

do_clean=false

while getopts c option
do
  case "${option}" in
    c) do_clean=true;;
    *) echo "unknown option"; exit 1;;
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
