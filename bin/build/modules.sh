#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

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
