#!/bin/bash

## Builds and starts the app once (no live reload).
##
## Usage:
##   ./bin/start.sh
##
## Environment:
##   - projectforge_encryption_key=TEMP_SECRET_KEY
##   - Loads variables from ./.env if present
##   - Sources $HOME/bin/oauth if present
##
## Requires:
##   - Go toolchain
##   - make
##
## Notes:
##   - Runs `make clean` and `make build`.
##   - Starts `./build/debug/projectforge -v --addr=0.0.0.0 all projectforge`.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"/..

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd make "install Xcode Command Line Tools or build-essential"

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"
export projectforge_encryption_key=TEMP_SECRET_KEY

# include env file
if [ -f ".env" ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    if [[ -n "$line" && ! $line =~ ^#.* ]]; then
      export "${line?}"
    fi
  done < ".env"
fi

make clean
make build
./build/debug/projectforge -v --addr=0.0.0.0 all projectforge
