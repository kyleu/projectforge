#!/bin/bash

## Builds a linux/amd64 binary, then builds and runs a Docker image.
##
## Usage:
##   ./bin/build/build-and-run-docker.sh
##
## Requires:
##   - Docker
##   - Go toolchain and make
##
## Notes:
##   - Temporarily moves ./build/debug/projectforge to repo root.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd docker "install Docker Desktop from https://www.docker.com/products/docker-desktop/"
require_cmd go "install Go from https://go.dev/dl/"
require_cmd make "install Xcode Command Line Tools or build-essential"

echo "Building [linux amd64]..."
GOOS=linux GOARCH=amd64 make build
mv ./build/debug/projectforge .
docker build -t=projectforge -f=./tools/release/Dockerfile.release .
rm ./projectforge
echo "Built Docker image [projectforge]"
docker run -it projectforge
