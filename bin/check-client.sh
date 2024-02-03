#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs eslint for the TypeScript project

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../client"

echo "=== linting client ==="
eslint --ext .js,.ts,.tsx .
