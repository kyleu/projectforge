#!/bin/bash

## Starts the app, without reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"/..

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
