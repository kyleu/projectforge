#!/bin/bash

## Packages each module directory into a zip archive.
##
## Usage:
##   ./bin/build/modules.sh
##
## Requires:
##   - zip
##
## Outputs:
##   - build/dist/projectforge_module_<module>.zip

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd zip "install zip from your package manager"

pwd

echo "packaging modules..."

function copyModule {
  echo "updating [$1] module"
  cd "$1"
  touch -- *
  touch .*
  zip -q -r -X "../../build/dist/projectforge_module_$1.zip" .
  cd ..
}

mkdir -p build/dist
cd module
for d in * ; do
  if [ -d "$d" ]; then
    copyModule "$d"
  fi
done
cd "$dir/../.."
echo "Artifacts written to ./build/dist/projectforge_module_*.zip"
