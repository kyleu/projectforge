#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

# $PF_SECTION_START(keys)$
# $PF_SECTION_END(keys)$

if command -v title &> /dev/null; then
  title "projectforge"
fi

[[ -f "$HOME/bin/oauth" ]] && . "$HOME/bin/oauth"
export projectforge_encryption_key=TEMP_SECRET_KEY

# include env file
if [ -f ".env" ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    if [[ ! $line =~ ^#.* ]]; then
      export "${line?}"
    fi
  done < ".env"
fi

./bin/templates.sh
go mod tidy

ulimit -S -n 65536

air
