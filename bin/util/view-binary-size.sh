#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

## Visualizes space usage for the release binary

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

make build-release
go tool nm -size build/release/projectforge | c++filt > ./tmp/nm.txt
