#!/bin/bash

## Builds the Docker image and runs it

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

echo "Building [linux amd64]..."
GOOS=linux GOARCH=amd64 make build
mv ./build/debug/{{{ .Exec }}} .
docker build -t={{{ .Key }}} -f=./tools/release/Dockerfile .
rm ./{{{ .Exec }}}
docker run -it {{{ .Key }}}
