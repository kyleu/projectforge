#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs goreleaser

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth

goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
