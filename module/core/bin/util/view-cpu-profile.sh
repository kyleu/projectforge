#!/bin/bash

## Serves a CPU profile with pprof on http://localhost:8000.
##
## Usage:
##   ./bin/util/view-cpu-profile.sh
##
## Requires:
##   - Go toolchain
##   - ./tmp/cpu.pprof generated beforehand

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"

echo "=== launching profiler for cpu.pprof ==="
go tool pprof -http=":8000" ./build/debug/{{{ .Exec }}} ./tmp/cpu.pprof
