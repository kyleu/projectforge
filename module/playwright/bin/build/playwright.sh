#!/bin/bash

## Runs the Playwright test suite.
##
## Usage:
##   ./bin/build/playwright.sh
##
## Requires:
##   - Node.js and npm
##   - Playwright dependencies installed in ./test/playwright

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."
cd "test/playwright"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd npx "install Node.js from https://nodejs.org"

echo "Testing application..."
npx playwright test

echo "Done!"
