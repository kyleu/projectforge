#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

if [ "$XSKIP_ANDROID" != "true" ]
then
  echo "building gomobile for Android..."
  mkdir -p build/dist/mobile_android_arm64
  time gomobile bind -o build/dist/mobile_android_arm64/$PF_EXECUTABLE$.aar -target=android $PF_PACKAGE$/app/cmd
  echo "gomobile for Android completed successfully, building distribution..."
  cd "build/dist/mobile_android_arm64"
  zip -r "../$PF_KEY$_${TGT}_mobile_android.zip" .
fi
