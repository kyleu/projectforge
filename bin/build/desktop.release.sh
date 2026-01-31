#!/bin/bash

## Builds desktop release artifacts for macOS, Linux, and Windows.
##
## Usage:
##   ./bin/build/desktop.release.sh [version] [-y|--yes]
##
## Arguments:
##   version  Version tag for output filenames (default: v0.0.0).
##   -y, --yes  Skip the confirmation prompt.
##
## Requires:
##   - Docker
##   - appdmg, codesign, lipo (macOS)
##   - APPLE_SIGNING_IDENTITY set for codesign
##   - curl and zip
##
## Outputs:
##   - build/dist/projectforge_<version>_darwin_*_desktop.(dmg|zip)
##   - build/dist/projectforge_<version>_linux_amd64_desktop.zip
##   - build/dist/projectforge_<version>_windows_amd64_desktop.zip

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/../.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_env() {
  if [[ -z "${!1:-}" ]]; then
    echo "error: required environment variable '$1' is not set" >&2
    exit 1
  fi
}

require_cmd docker "install Docker Desktop from https://www.docker.com/products/docker-desktop/"
require_cmd zip "install zip from your package manager"
require_cmd curl "install curl from your package manager"
require_cmd appdmg "install via npm i -g appdmg"
require_cmd codesign "requires macOS codesign tool"
require_cmd lipo "requires macOS Xcode command line tools"
require_env APPLE_SIGNING_IDENTITY

auto_yes=false
TGT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -y|--yes) auto_yes=true; shift;;
    --) shift; break;;
    -*)
      echo "unknown option: $1" >&2
      exit 1
      ;;
    *)
      if [[ -z "$TGT" ]]; then
        TGT="$1"
        shift
      else
        echo "unexpected argument: $1" >&2
        exit 1
      fi
      ;;
  esac
done

TGT=${TGT:-v0.0.0}

if ! $auto_yes; then
  read -r -p "Build desktop release artifacts for ${TGT}? [y/N] " confirm
  case "$confirm" in
    [yY][eE][sS]|[yY]) ;;
    *) echo "aborted"; exit 0;;
  esac
fi

if command -v retry &> /dev/null
then
  retry -t 4 -- docker build -f tools/desktop/Dockerfile.desktop -t projectforge .
else
  docker build -f tools/desktop/Dockerfile.desktop -t projectforge .
fi


rm -rf tmp/release
mkdir -p tmp/release

cd "tmp/release"

id=$(docker create projectforge)
docker cp "$id":/dist - > ./desktop.tar
docker rm -v "$id"
tar -xvf "desktop.tar"
rm "desktop.tar"

mv dist/darwin_amd64/projectforge "projectforge.darwin"
mv dist/darwin_arm64/projectforge "projectforge.darwin.arm64"
mv dist/linux_amd64/projectforge "projectforge"
mv dist/windows_amd64/projectforge "projectforge.exe"
rm -rf "dist"

# darwin amd64
cp -R "../../tools/desktop/template" .

mkdir -p "./Project Forge.app/Contents/Resources"
mkdir -p "./Project Forge.app/Contents/MacOS"

cp -R "./template/darwin/Info.plist" "./Project Forge.app/Contents/Info.plist"
cp -R "./template/darwin/icons.icns" "./Project Forge.app/Contents/Resources/icons.icns"

cp "projectforge.darwin" "./Project Forge.app/Contents/MacOS/projectforge"

echo "signing amd64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app/Contents/MacOS/projectforge"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app"

cp "./template/darwin/appdmg.config.json" "./appdmg.config.json"

echo "building macOS amd64 DMG..."
appdmg "appdmg.config.json" "./projectforge_${TGT}_darwin_amd64_desktop.dmg"
zip -r "projectforge_${TGT}_darwin_amd64_desktop.zip" "./Project Forge.app"

# darwin arm64
cp "projectforge.darwin.arm64" "./Project Forge.app/Contents/MacOS/projectforge"

echo "signing arm64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app/Contents/MacOS/projectforge"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app"

echo "building macOS arm64 DMG..."
appdmg "appdmg.config.json" "./projectforge_${TGT}_darwin_arm64_desktop.dmg"
zip -r "projectforge_${TGT}_darwin_arm64_desktop.zip" "./Project Forge.app"

# macOS universal
rm "./Project Forge.app/Contents/MacOS/projectforge"
lipo -create -output "./Project Forge.app/Contents/MacOS/projectforge" projectforge.darwin projectforge.darwin.arm64

echo "signing universal desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app/Contents/MacOS/projectforge"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s "${APPLE_SIGNING_IDENTITY}" "./Project Forge.app"

echo "building macOS universal DMG..."
appdmg "appdmg.config.json" "./projectforge_${TGT}_darwin_all_desktop.dmg"
zip -r "projectforge_${TGT}_darwin_all_desktop.zip" "./Project Forge.app"

# linux
echo "building Linux zip..."
zip "projectforge_${TGT}_linux_amd64_desktop.zip" "./projectforge"

#windows
echo "building Windows zip..."
curl -L -o webview.dll https://github.com/webview/webview/raw/master/dll/x64/webview.dll
curl -L -o WebView2Loader.dll https://github.com/webview/webview/raw/master/dll/x64/WebView2Loader.dll
zip "projectforge_${TGT}_windows_amd64_desktop.zip" "./projectforge.exe" "./webview.dll" "./WebView2Loader.dll"

mkdir -p "../../build/dist"
mv "./projectforge_${TGT}_darwin_amd64_desktop.dmg" "../../build/dist"
mv "./projectforge_${TGT}_darwin_amd64_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_darwin_arm64_desktop.dmg" "../../build/dist"
mv "./projectforge_${TGT}_darwin_arm64_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_darwin_all_desktop.dmg" "../../build/dist"
mv "./projectforge_${TGT}_darwin_all_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_linux_amd64_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_windows_amd64_desktop.zip" "../../build/dist"

cd "$dir/../.."
echo "Builds written to ./build/dist (projectforge_${TGT}_*_desktop.*)"
