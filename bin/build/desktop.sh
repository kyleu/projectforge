#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Uses `tools/desktop` to build a desktop application

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

cd tools/desktop

go mod tidy

make build
echo "build complete, starting desktop application"

cd $dir/../..
build/debug/projectforge-desktop
