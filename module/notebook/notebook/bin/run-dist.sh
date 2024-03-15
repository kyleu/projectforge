#!/bin/bash

## Builds the notebook and runs it

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

./bin/build.sh
npx http-server ./dist
