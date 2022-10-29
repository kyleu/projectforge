#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the TypeScript resources, then watches for changes via `watchexec`
## Requires node and tsc available on the path

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../../client

echo "Watching TypeScript compilation for [client/src]..."
watchexec --exts css,ts,tsx node build.js
