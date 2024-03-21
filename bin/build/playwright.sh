#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs the Playwright tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

echo "Testing application..."
npm playwright test

echo "Done!"
