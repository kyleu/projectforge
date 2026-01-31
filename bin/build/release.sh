#!/bin/bash

## Runs goreleaser for an official release.
##
## Usage:
##   ./bin/build/release.sh [-y|--yes]
##
## Arguments:
##   -y, --yes  Skip the confirmation prompt.
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

auto_yes=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    -y|--yes) auto_yes=true; shift;;
    --) shift; break;;
    *) echo "unknown option: $1" >&2; exit 1;;
  esac
done

if ! $auto_yes; then
  read -r -p "Run goreleaser release? [y/N] " confirm
  case "$confirm" in
    [yY][eE][sS]|[yY]) ;;
    *) echo "aborted"; exit 0;;
  esac
fi

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"

goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
echo "Output written to [./dist] (snapshot)"
