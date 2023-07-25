#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Builds the app (or just use make build)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

echo "Building publishing Docker image..."
docker build -f tools/publish/Dockerfile.publish -t projectforge-publish .

# TODO: Remove
docker inspect -f "Size: {{ .Size }}" projectforge-publish
docker run --platform=linux/amd64 -it projectforge-publish
title dockertest
