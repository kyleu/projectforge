#!/bin/bash

## Runs goreleaser in test mode

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

./bin/asset-embed.sh
goreleaser -f ./tools/release/.goreleaser.yml --snapshot --skip-publish --rm-dist
./bin/asset-reset.sh
