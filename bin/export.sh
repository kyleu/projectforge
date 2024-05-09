#!/bin/bash
## Generates a single project

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

act=${1:-preview}
${dir}/../bin/build/build.sh darwin arm64 && ${dir}/../build/darwin/arm64/projectforge $act
