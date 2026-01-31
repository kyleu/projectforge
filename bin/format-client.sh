#!/bin/bash

## Formats the client TypeScript project.
##
## Usage:
##   ./bin/format-client.sh
##
## Requires:
##   - Node.js and npm
##   - Client deps installed (npm install in ./client)
##
## Notes:
##   - Runs `npm run format`.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."
cd "client"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd npm "install Node.js from https://nodejs.org"

echo "=== formatting ==="
npm run format
