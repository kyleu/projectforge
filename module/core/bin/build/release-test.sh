#!/bin/bash

## Runs goreleaser in snapshot (test) mode.
##
## Usage:
##   ./bin/build/release-test.sh
##
##
## Environment:
##   - PUBLISH_TEST=true (set by this script)
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

if grep -q "^dockers_v2:" ./tools/release/.goreleaser.yml; then
  require_cmd docker "install from https://docs.docker.com/get-docker/"

  export BUILDX_BUILDER="${BUILDX_BUILDER:-{{{ .Key }}}-release}"
  if ! docker buildx inspect "$BUILDX_BUILDER" >/dev/null 2>&1; then
    docker buildx create --name "$BUILDX_BUILDER" --driver docker-container --use >/dev/null
  else
    docker buildx use "$BUILDX_BUILDER" >/dev/null
  fi
  docker buildx inspect --bootstrap "$BUILDX_BUILDER" >/dev/null
fi

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"

export PUBLISH_TEST=true
goreleaser -f ./tools/release/.goreleaser.yml --snapshot --timeout 240m --clean --verbose
echo "Output written to [./dist] (snapshot)"
unset PUBLISH_TEST
