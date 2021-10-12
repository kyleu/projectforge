#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

docker build -f tools/desktop/Dockerfile -t projectforge .

rm -rf tmp/release
mkdir -p tmp/release

cd "tmp/release"

id=$(docker create projectforge)
docker cp $id:/dist - > ./desktop.tar
docker rm -v $id
tar -xvf "desktop.tar"
rm "desktop.tar"

mv dist/darwin_amd64/projectforge "projectforge.macos"
mv dist/linux_amd64/projectforge "projectforge"
mv dist/windows_amd64/projectforge "projectforge.exe"
rm -rf "dist"

# macOS
cp -R "../../tools/desktop/template" .

mkdir -p "./Project Forge.app/Contents/Resources"
mkdir -p "./Project Forge.app/Contents/MacOS"

cp -R "./template/macos/Info.plist" "./Project Forge.app/Contents/Info.plist"
cp -R "./template/macOS/icons.icns" "./Project Forge.app/Contents/Resources/icons.icns"

cp "projectforge.macos" "./Project Forge.app/Contents/MacOS/projectforge"

echo "signing desktop binary..."
codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./Project Forge.app/Contents/MacOS/projectforge"
codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./Project Forge.app"

cp "./template/macos/appdmg.config.json" "./appdmg.config.json"

echo "building macOS DMG..."
appdmg "appdmg.config.json" "./projectforge_${TGT}_macos_x86_64_desktop.dmg"
zip -r "projectforge_${TGT}_macos_x86_64_desktop.zip" "./Project Forge.app"

echo "building Linux zip..."
zip "projectforge_${TGT}_linux_x86_64_desktop.zip" "./projectforge"

echo "building Windows zip..."
curl -L -o webview.dll https://github.com/webview/webview/raw/master/dll/x64/webview.dll
curl -L -o WebView2Loader.dll https://github.com/webview/webview/raw/master/dll/x64/WebView2Loader.dll
zip "projectforge_${TGT}_windows_x86_64_desktop.zip" "./projectforge.exe" "./webview.dll" "./WebView2Loader.dll"

mkdir -p "../../build/dist"
mv "./projectforge_${TGT}_macos_x86_64_desktop.dmg" "../../build/dist"
mv "./projectforge_${TGT}_macos_x86_64_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_linux_x86_64_desktop.zip" "../../build/dist"
mv "./projectforge_${TGT}_windows_x86_64_desktop.zip" "../../build/dist"
