#!/bin/bash

## Compiles quicktemplate templates in ./views, skipping unchanged inputs.
##
## Usage:
##   ./bin/templates.sh [force]
##
## Arguments:
##   force  Rebuild even if hashes match.
##
## Requires:
##   - qtc in PATH
##
## Notes:
##   - Writes hashes to ./tmp/<dir>.hashcode to detect changes.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

FORCE="${1:-}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd qtc "install via 'go install github.com/valyala/quicktemplate/qtc@latest'"

hash_cmd="md5sum"
if ! command -v "$hash_cmd" >/dev/null 2>&1; then
  hash_cmd="md5"
fi
if ! command -v "$hash_cmd" >/dev/null 2>&1; then
  echo "error: required command 'md5sum' or 'md5' not found" >&2
  exit 1
fi

function tmpl {
  echo "updating [$1] templates"
  if test -f "$ftgt"; then
    mv "$ftgt" "$fsrc"
  fi
  qtc -ext "$2" -dir "$1" 2> >(grep -v Compiling)
}

function check {
  if [[ -d "$1" ]]; then
    fsrc="tmp/$1.hashcode"
    ftgt="tmp/$1.hashcode.tmp"

    mkdir -p tmp/

    find "$1" -type f | grep "\.$2$" | xargs "$hash_cmd" > "$ftgt"

    if cmp -s "$fsrc" "$ftgt"; then
      if [ "$FORCE" = "force" ]; then
        tmpl "$1" "$2"
      else
        rm "$ftgt"
      fi
    else
      tmpl "$1" "$2"
    fi
  fi
}

check "views" "html"
