#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

pwd

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

echo "packaging modules..."
mkdir -p build/dist/module

function z {
  echo "updating [$1] module"
  cd module/$1
  zip -r "../../build/dist/projectforge_module_$1.zip" .
  cd ../..
}

z core
z database
z desktop
z marketing
z migration
z mobile
z oauth
z postgres
z search
z sqlite
