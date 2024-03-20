#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Runs goreleaser in test mode

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth

export PUBLISH_TEST=true
goreleaser -f ./tools/release/.goreleaser.yml --snapshot --timeout 240m --skip-publish --clean --verbose
unset PUBLISH_TEST
