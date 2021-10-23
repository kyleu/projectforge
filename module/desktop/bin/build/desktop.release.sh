#!/bin/bash

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

mv dist/darwin_amd64/{{{ .Key }}} "{{{ .Key }}}.macos"
mv dist/darwin_arm64/{{{ .Key }}} "{{{ .Key }}}.macos.arm64"
mv dist/linux_amd64/{{{ .Key }}} "{{{ .Key }}}"
mv dist/windows_amd64/{{{ .Key }}} "{{{ .Key }}}.exe"
rm -rf "dist"

# macOS x86_64
cp -R "../../tools/desktop/template" .

mkdir -p "./Project Forge.app/Contents/Resources"
mkdir -p "./Project Forge.app/Contents/MacOS"

cp -R "./template/macos/Info.plist" "./Project Forge.app/Contents/Info.plist"
cp -R "./template/macOS/icons.icns" "./Project Forge.app/Contents/Resources/icons.icns"

cp "{{{ .Key }}}.macos" "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"

echo "signing amd64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app"

cp "./template/macos/appdmg.config.json" "./appdmg.config.json"

echo "building macOS amd64 DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_macos_x86_64_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_macos_x86_64_desktop.zip" "./Project Forge.app"

# macOS arm64
cp "{{{ .Key }}}.macos.arm64" "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"

echo "signing arm64 desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app"

echo "building macOS arm64 DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_macos_arm64_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_macos_arm64_desktop.zip" "./Project Forge.app"

# macOS universal
rm "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"
lipo -create -output "./Project Forge.app/Contents/MacOS/{{{ .Key }}}" {{{ .Key }}}.macos {{{ .Key }}}.macos.arm64

echo "signing universal desktop binary..."
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app/Contents/MacOS/{{{ .Key }}}"
codesign -f --options=runtime --verbose=4 --deep --force --strict -s '{{{ .Info.SigningIdentity }}}' "./Project Forge.app"

echo "building macOS universal DMG..."
appdmg "appdmg.config.json" "./{{{ .Key }}}_${TGT}_macos_all_desktop.dmg"
zip -r "{{{ .Key }}}_${TGT}_macos_all_desktop.zip" "./Project Forge.app"

# linux
echo "building Linux zip..."
zip "{{{ .Key }}}_${TGT}_linux_x86_64_desktop.zip" "./{{{ .Key }}}"

#windows
echo "building Windows zip..."
curl -L -o webview.dll https://github.com/webview/webview/raw/master/dll/x64/webview.dll
curl -L -o WebView2Loader.dll https://github.com/webview/webview/raw/master/dll/x64/WebView2Loader.dll
zip "{{{ .Key }}}_${TGT}_windows_x86_64_desktop.zip" "./{{{ .Key }}}.exe" "./webview.dll" "./WebView2Loader.dll"

mkdir -p "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_x86_64_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_x86_64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_arm64_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_arm64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_all_desktop.dmg" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_macos_all_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_linux_x86_64_desktop.zip" "../../build/dist"
mv "./{{{ .Key }}}_${TGT}_windows_x86_64_desktop.zip" "../../build/dist"
