#!/bin/bash

## Builds Android AAR and debug APK zips.
##
## Usage:
##   ./bin/build/android.sh [version]
##
## Arguments:
##   version  Version tag for output filenames (default: 0.0.0).
##
## Requires:
##   - Go toolchain and gomobile
##   - Android SDK/NDK, Gradle, and JDK
##   - zip
##
## Outputs:
##   - build/dist/projectforge_<version>_android_aar.zip
##   - build/dist/projectforge_<version>_android_apk.zip

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
require_cmd gradle "install Gradle or use a wrapper"
require_cmd zip "install zip from your package manager"

TGT=$1
[ "$TGT" ] || TGT="0.0.0"

echo "building gomobile for Android..."
mkdir -p build/dist/mobile_android_arm64
GOARCH=arm64 time gomobile bind -o build/dist/mobile_android_arm64/projectforge.aar -target=android -androidapi 21 projectforge.dev/projectforge/app/cmd
echo "gomobile for Android completed successfully, building distribution..."
cd "build/dist/mobile_android_arm64"
zip -r "../projectforge_${TGT}_android_aar.zip" .

echo "creating Android project..."
cd "$dir/../.."
mkdir -p build/dist/mobile_android_app_arm64
cp -R tools/android/* build/dist/mobile_android_app_arm64

echo "building Android project..."
cd build/dist/mobile_android_app_arm64
rm -rf ./app/libs
mkdir -p ./app/libs
cp ../mobile_android_arm64/projectforge.aar ./app/libs/
cp ../mobile_android_arm64/projectforge-sources.jar ./app/libs/

gradle assembleDebug
cd app/build/outputs/apk/debug
zip -r "$dir/../../build/dist/projectforge_${TGT}_android_apk.zip" .

echo "Output written to [./build/dist/projectforge_${TGT}_android_aar.zip] and [./build/dist/projectforge_${TGT}_android_apk.zip]"
