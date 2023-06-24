#!/bin/bash

## Meant to be run as part of the release process, builds desktop apps

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

docker build -f tools/desktop/Dockerfile -t {{{ .Key }}} .

rm -rf tmp/release
mkdir -p tmp/release

cd "tmp/release"

id=$(docker create {{{ .Key }}})
docker cp $id:/dist - > ./desktop.tar
docker rm -v $id
tar -xvf "desktop.tar"
rm "desktop.tar"

mv dist/darwin_amd64/{{{ .Key }}} "{{{ .Key }}}.darwin"
mv dist/darwin_arm64/{{{ .Key }}} "{{{ .Key }}}.darwin.arm64"
mv dist/linux_amd64/{{{ .Key }}} "{{{ .Key }}}"
mv dist/windows_amd64/{{{ .Key }}} "{{{ .Key }}}.exe"
rm -rf "dist"

# darwin amd64
cp -R "../../tools/desktop/template" .

mkdir -p "./{{{ .Name }}}.app/Contents/Resources"
mkdir -p "./{{{ .Name }}}.app/Contents/MacOS"

cp -R "./template/macos/Info.plist" "./{{{ .Name }}}.app/Contents/Info.plist"
cp -R "./template/macos/icons.icns" "./{{{ .Name }}}.app/Contents/Resources/icons.icns"

cp "{{{ .Key }}}.macos" "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"

echo "signing amd64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app"

cp "./template/macos/appdmg.config.json" "./appdmg.config.json"

echo "building macOS amd64 DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_darwin_amd64_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_darwin_amd64_desktop.zip" "./{{{ .Name }}}.app"

# darwin arm64
cp "{{{ .Key }}}.darwin.arm64" "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"

echo "signing arm64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app"

echo "building macOS arm64 DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_darwin_arm64_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_darwin_arm64_desktop.zip" "./{{{ .Name }}}.app"

# macOS universal
rm "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"
lipo -create -output "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}" {{{ .Key }}}.darwin {{{ .Key }}}.darwin.arm64

echo "signing universal desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./{{{ .Name }}}.app"

echo "building macOS universal DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_darwin_all_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_darwin_all_desktop.zip" "./{{{ .Name }}}.app"

# linux
echo "building Linux zip..."
zip "{{{ .Key }}}_${TGT}_linux_amd64_desktop.zip" "./{{{ .Key }}}"

#windows
echo "building Windows zip..."
curl -L -o webview.dll https://github.com/webview/webview/raw/master/dll/x64/webview.dll
curl -L -o WebView2Loader.dll https://github.com/webview/webview/raw/master/dll/x64/WebView2Loader.dll
zip "{{{ .Key }}}_${TGT}_windows_amd64_desktop.zip" "./{{{ .Key }}}.exe" "./webview.dll" "./WebView2Loader.dll"

mkdir -p "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_amd64_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_amd64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_arm64_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_arm64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_all_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_darwin_all_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_linux_amd64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_windows_x86_64_desktop.zip" "../../build/dist"
