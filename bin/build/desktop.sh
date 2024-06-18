#!/bin/bash

## Uses `tools/desktop` to build a desktop application

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

cd tools/desktop

go mod tidy

make build
echo "build complete, starting desktop application"

cd "$dir/../.."
build/debug/projectforge-desktop
