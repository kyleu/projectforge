#!/bin/bash

## Embeds assets for building into the project, only need for release builds. Run `asset-reset.sh` to undo before committing

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

echo "Embedding assets..."
go-embed -input assets -output app/assets/assets.go
git update-index --assume-unchanged app/assets/assets.go