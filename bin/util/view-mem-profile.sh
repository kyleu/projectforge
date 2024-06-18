#!/bin/bash

## Starts a pprof server using the (previously-recorded) heap dump at ./tmp/mem.pprof

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

echo "=== launching profiler for mem.pprof ==="
go tool pprof -http=":8000" ./build/debug/projectforge ./tmp/mem.pprof
