#!/bin/bash

## Builds all the templates using quicktemplate, skipping if unchanged

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

FORCE="${1:-}"

function tmpl {
  echo "updating [$1] templates"
  if test -f "$ftgt"; then
    mv "$ftgt" "$fsrc"
  fi
  qtc -ext $2 -dir "$1" 2> >(grep -v Compiling)
}

function check {
  if [[ -d "$1" ]]; then
    fsrc="tmp/$1.hashcode"
    ftgt="tmp/$1.hashcode.tmp"

    mkdir -p tmp/

    m="md5sum"
    if ! [ -x "$(command -v $m)" ]; then
      m="md5"
    fi

    find "$1" -type f | grep "\.$2$" | xargs $m > "$ftgt"

    if cmp -s "$fsrc" "$ftgt"; then
      if [ "$FORCE" = "force" ]; then
        tmpl $1 $2
      else
        rm "$ftgt"
      fi
    else
      tmpl $1 $2
    fi
  fi
}

check "views" "html"
