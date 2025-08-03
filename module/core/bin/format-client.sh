#!/bin/bash

## Formatting code from the TypeScript projects

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."
cd "client"

echo "=== formatting ==="
npm run format
