#!/bin/bash
# $PF_GENERATE_ONCE$
## Creates symlinks in the data directory

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

create_symlink() {
    local target_directory="$1"
    local link_directory="$2"
    if [ ! -e "$link_directory" ]; then
        ln -s "$target_directory" "$link_directory"
        echo "symlink [$link_directory -> $target_directory] created"
    fi
}

mkdir -p "docs/data"
cd "docs/data"

create_symlink "../../../{{{ .Key }}}.sqlite" "{{{ .Key }}}.sqlite"
