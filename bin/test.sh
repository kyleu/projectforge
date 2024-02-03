#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs all the tests. Pass "-c" to clear the cache first, "-w" to watch for changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

do_clean=false
do_watch=false
target_test=""

while getopts cw option
do
  case "${option}" in
    c) do_clean=true;;
    w) do_watch=true;;
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
	export $(cat test.env | grep -v "#" | xargs)
fi

if [ -f "./bin/test-setup.sh" ]; then
	./bin/test-setup.sh
fi

gotestsum${test_args} -- -race ./app/...
