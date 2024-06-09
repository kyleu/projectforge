#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Links the application support directory for the current user to ./cfg

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

tgt="$HOME/Library/Application Support/Project Forge"

mkdir -p "$tgt"
ln -s "$tgt" "cfg"
echo "configuration directory linked"
