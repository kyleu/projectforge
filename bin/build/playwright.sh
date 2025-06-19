#!/bin/bash

## Runs the Playwright tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."
cd "test/playwright"

echo "Testing application..."
npx playwright test

echo "Done!"
