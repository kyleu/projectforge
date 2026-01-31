#!/bin/bash

## Builds the iOS xcframework and app zip.
##
## Usage:
##   ./bin/build/ios.sh [version]
##
## Arguments:
##   version  Version tag for output filenames (default: 0.0.0).
##
## Requires:
##   - Go toolchain and gomobile
##   - Xcode (xcodebuild) and xcodegen
##   - zip
##
## Outputs:
##   - build/dist/projectforge_<version>_ios_framework.zip
##   - build/dist/projectforge_<version>_ios_app.zip

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd go "install Go from https://go.dev/dl/"
require_cmd gomobile "install via 'go install golang.org/x/mobile/cmd/gomobile@latest'"
require_cmd xcodebuild "install Xcode from the App Store"
require_cmd xcodegen "install via brew install xcodegen"
require_cmd zip "install zip from your package manager"

TGT=$1
[ "$TGT" ] || TGT="0.0.0"

echo "building gomobile for iOS..."
GOARCH=arm64 time gomobile bind -o build/dist/mobile_ios_arm64/projectforgeServer.xcframework -target=ios projectforge.dev/projectforge/app/cmd
echo "gomobile for iOS completed successfully, building distribution..."
cd "build/dist/mobile_ios_arm64/projectforgeServer.xcframework"
zip --symlinks -r "../../projectforge_${TGT}_ios_framework.zip" .

echo "Building iOS app..."
cd "$dir/../../tools/ios"

rm -rf ../../build/dist/mobile_ios_app_arm64
mkdir -p ../../build/dist/mobile_ios_app_arm64

xcodegen generate --spec xcodegen.yml --project ../../build/dist/mobile_ios_app_arm64

mv Info.plist ../../build/dist/mobile_ios_app_arm64
cd ../../build/dist/mobile_ios_app_arm64

xcodebuild -project "Project Forge.xcodeproj" -allowProvisioningUpdates
zip -r "$dir/../../build/dist/projectforge_${TGT}_ios_app.zip" .

cd "$dir/../.."
echo "Output written to ./build/dist/projectforge_${TGT}_ios_framework.zip and ./build/dist/projectforge_${TGT}_ios_app.zip"
