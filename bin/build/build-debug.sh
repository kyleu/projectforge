#!/bin/bash

## Builds the app in debug mode (or just use "make build" with an appropriate GOOS and GOARCH)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

os=${1:-darwin}
arch=${2:-amd64}
fn=${3:-projectforge}

echo "Building debug $fn [$os $arch]..."
env GOOS="$os" GOARCH="$arch" make build
mkdir -p "./build/$os/$arch"
mv "./build/debug/$fn" "./build/$os/$arch/$fn"
