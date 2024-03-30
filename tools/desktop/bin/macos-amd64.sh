#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Uses `tools/desktop` to build a desktop application, intended to be run from Docker

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

echo "starting macOS amd64 desktop build..."
GODEBUG=asyncpreemptoff=1 GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ go build -o ../../dist/darwin_amd64/projectforge
