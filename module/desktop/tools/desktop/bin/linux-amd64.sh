#!/bin/bash

## Uses `tools/desktop` to build a desktop application, intended to be run from Docker

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

echo "starting Linux amd64 desktop build..."
GOOS=linux GOARCH=amd64 go build -o ../../dist/linux_amd64/{{{ .Exec }}}
