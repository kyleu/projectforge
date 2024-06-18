#!/bin/bash

## Uses gomod to visualize the module graph
## Requires gomod available on the path

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

echo "building dependency SVG..."
gomod graph | dot -Tsvg -o ./tmp/deps.svg

open ./tmp/deps.svg
