#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

if [ "$XSKIP_IOS" != "true" ]
then
  echo "building gomobile for iOS..."
  time gomobile bind -o build/dist/mobile_ios_arm64/{{{ .Exec }}}.framework -target=ios {{{ .Package }}}/app/cmd
  echo "gomobile for iOS completed successfully, building distribution..."
  cd "build/dist/mobile_ios_arm64/{{{ .Key }}}.framework"
  zip --symlinks -r "../../{{{ .Key }}}_${TGT}_mobile_ios.zip" .
fi
