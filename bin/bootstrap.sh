#!/bin/bash

## Installs base Go tooling used by Project Forge scripts.
##
## Usage:
##   ./bin/bootstrap.sh
##
## Requires:
##   - Go toolchain in PATH
##
## Installs:
##   - github.com/air-verse/air@latest
##   - github.com/valyala/quicktemplate/qtc@latest
##   - gotest.tools/gotestsum@latest
##   - mvdan.cc/gofumpt@latest
##   - github.com/nikolaydubina/go-cover-treemap@latest
##   - github.com/go-delve/delve/cmd/dlv@latest
##   - golang.org/x/mobile/cmd/gomobile@latest
##   - golang.org/x/mobile/cmd/gobind@latest
##
## Notes:
##   - Binaries are installed to GOPATH/bin or GOBIN
##   - Ensure GOPATH\bin or GOBIN is on PATH

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"

go install github.com/air-verse/air@latest
go install github.com/valyala/quicktemplate/qtc@latest
go install gotest.tools/gotestsum@latest
go install mvdan.cc/gofumpt@latest
go install github.com/nikolaydubina/go-cover-treemap@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
