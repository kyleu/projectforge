#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the app, installing all prerequisites

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

echo "Installing utilities..."
./bin/bootstrap.sh

echo "Compiling templates..."
./bin/templates.sh

echo "Downloading dependencies..."
go mod download

echo "Updating client..."
cd client
npm install
cd ..

echo "Building client..."
./bin/build/client.sh

echo "Building application..."
make build-release

echo "Done!"
