#!/bin/bash

## Builds the app in release mode (or just use "make build-release" with an appropriate GOOS and GOARCH)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

os=${1:-darwin}
arch=${2:-amd64}
fn=${3:-projectforge}

echo "Building $fn [$os $arch]..."
env GOOS="$os" GOARCH="$arch" make build-release
mkdir -p "./build/$os/$arch"
mv "./build/release/$fn" "./build/$os/$arch/$fn"
