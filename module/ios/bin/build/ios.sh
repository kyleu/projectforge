#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

echo "building gomobile for iOS..."
time gomobile bind -o build/dist/mobile_ios_arm64/{{{ .Key }}}.framework -target=ios {{{ .Package }}}/app/cmd
echo "gomobile for iOS completed successfully, building distribution..."
cd "build/dist/mobile_ios_arm64/{{{ .Key }}}.framework"
zip --symlinks -r "../../{{{ .Key }}}_${TGT}_mobile_ios_framework.zip" .

echo "Building iOS app..."
cd $dir/../../tools/ios

rm -rf {{{ .Key }}}.framework
cp -R ../../build/dist/mobile_ios_arm64/{{{ .Key }}}.framework ./{{{ .Key }}}.framework

xcodebuild -project app.xcodeproj

cd build/Release-iphoneos/

zip -r "$dir/../../build/dist/{{{ .Key }}}_${TGT}_mobile_ios_app.zip" "{{{ .Key }}}.app"
