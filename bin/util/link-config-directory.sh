#!/bin/bash

## Symlinks the macOS Application Support directory to ./cfg.
##
## Usage:
##   ./bin/util/link-config-directory.sh
##
## Requires:
##   - macOS (uses ~/Library/Application Support)
##
## Notes:
##   - Creates the target directory if needed.
##   - Fails if ./cfg already exists.

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

tgt="$HOME/Library/Application Support/Project Forge"

mkdir -p "$tgt"
ln -s "$tgt" "cfg"
echo "configuration directory linked"
