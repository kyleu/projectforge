#!/bin/bash

## Builds the app (or just use make build)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

echo "Building Docker \"publish\" image..."
docker build -f tools/publish/Dockerfile.publish -t {{{ .Exec }}}-publish .

# TODO: Remove
docker inspect -f "Size: {{ .Size }}" {{{ .Exec }}}-publish
docker run --platform=linux/amd64 -it {{{ .Exec }}}-publish
title dockertest
