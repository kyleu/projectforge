#!/bin/bash

## Runs the Playwright tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

echo "Testing application..."
npm playwright test

echo "Done!"
