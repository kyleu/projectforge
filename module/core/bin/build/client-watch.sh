#!/bin/bash

## Builds the TypeScript resources, then watches for changes via `watchexec`
## Requires node and tsc available on the path

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../../client"

echo "Watching TypeScript compilation for [client/src]..."
watchexec --exts css,ts,tsx --no-vcs-ignore node build.js
