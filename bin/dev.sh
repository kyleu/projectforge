#!/bin/bash

## Starts the dev server with live reload (air).
##
## Usage:
##   ./bin/dev.sh
##
## Environment:
##   - Loads variables from ./.env if present
##   - Sources $HOME/bin/oauth if present
##
## Requires:
##   - Go toolchain
##   - air (live reload)
##   - qtc (templates, via ./bin/templates.sh)
##
## Notes:
##   - Runs ./bin/templates.sh and `go mod tidy` before first run.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd air "install via 'go install github.com/air-verse/air@latest'"
require_cmd qtc "install via 'go install github.com/valyala/quicktemplate/qtc@latest'"

# $PF_SECTION_START(keys)$
# $PF_SECTION_END(keys)$

if command -v title &> /dev/null; then
  title "projectforge"
fi

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"
export projectforge_encryption_key=TEMP_SECRET_KEY
export GOEXPERIMENT=jsonv2

# include env file
if [ -f ".env" ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    if [[ -n "$line" && ! $line =~ ^#.* ]]; then
      export "${line?}"
    fi
  done < ".env"
fi

./bin/templates.sh
go mod tidy

air
