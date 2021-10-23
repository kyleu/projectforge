#!/bin/bash

## Uses `tools/desktop` to build a desktop application, intended to be run from Docker

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir

go mod tidy

echo "starting macOS desktop build..."
GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ go build -o ../../dist/darwin_amd64/projectforge

echo "starting macOS arm64 desktop build..."
GOOS=darwin GOARCH=arm64 CC=aarch64-apple-darwin20.2-clang CXX=aarch64-apple-darwin20.2-clang++ go build -o ../../dist/darwin_arm64/projectforge

echo "starting Linux desktop build..."
GOOS=linux GOARCH=amd64 go build -o ../../dist/linux_amd64/projectforge

echo "starting Windows desktop build..."
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o ../../dist/windows_amd64/projectforge
