#!/bin/bash

## Builds client assets and watches for changes.
##
## Usage:
##   ./bin/build/client-watch.sh
##
## Requires:
##   - Node.js and npm
##   - watchexec in PATH
##
## Notes:
##   - Watches *.css, *.ts, *.tsx and runs `node build.js` in ./client.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../../client"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd node "install Node.js from https://nodejs.org"
require_cmd npm "install Node.js from https://nodejs.org"
require_cmd watchexec "install from https://github.com/watchexec/watchexec"

echo "Watching TypeScript compilation for [client/src]..."
echo "Output written to [../assets] on each successful build"
watchexec --exts css,ts,tsx --no-vcs-ignore node build.js
