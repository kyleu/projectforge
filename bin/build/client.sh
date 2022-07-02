#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Uses `esbuild` to compile the scripts in `client`
## Requires node, tsc, and esbuild available on the path

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../../client

node build.js
