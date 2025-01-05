#!/bin/bash

## Runs eslint for the TypeScript project

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../client"

echo "=== linting client ==="
npm run lint
