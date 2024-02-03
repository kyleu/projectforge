#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the app (or just use make build)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

cd client
npm install
cd ..
./bin/build/client.sh

make clean
./bin/templates.sh
go mod tidy
./bin/build/build.sh
