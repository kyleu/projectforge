#!/bin/bash

## Runs goreleaser for an official release.
##
## Usage:
##   ./bin/build/release.sh
##
## Environment:
##   - Sources $HOME/bin/oauth if present
##
## Requires:
##   - goreleaser in PATH

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd goreleaser "install from https://goreleaser.com/install/"

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"

goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
echo "Output written to [./dist] (snapshot)"
