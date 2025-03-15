#!/bin/bash

## Starts the notebook in dev mode

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

if [ -f ".env" ]; then
  while IFS= read -r line || [ -n "$line" ]; do
    if [[ -n "$line" && ! $line =~ ^#.* ]]; then
      export "${line?}"
    fi
  done < ".env"
fi

npm run dev -- --port {{{ .NotebookPort }}}
