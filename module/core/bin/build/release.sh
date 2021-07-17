#!/bin/bash

## Runs goreleaser

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

./bin/asset-embed.sh
goreleaser -f ./tools/release/.goreleaser.yml release --rm-dist
./bin/asset-reset.sh
