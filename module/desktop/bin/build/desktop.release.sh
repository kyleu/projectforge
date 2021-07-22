#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

if [ "$XSKIP_DESKTOP" != "true" ]
then
  docker build -f tools/desktop/Dockerfile -t $PF_KEY$ .

  rm -rf tmp/release
  mkdir -p tmp/release

  cd "tmp/release"

  id=$(docker create $PF_KEY$)
  docker cp $id:/dist - > ./desktop.tar
  docker rm -v $id
  tar -xvf "desktop.tar"
  rm "desktop.tar"

  mv dist/darwin_amd64/$PF_KEY$ "$PF_KEY$.macos"
  mv dist/linux_amd64/$PF_KEY$ "$PF_KEY$"
  mv dist/windows_amd64/$PF_KEY$ "$PF_KEY$.exe"
  rm -rf "dist"

  # macOS
  cp -R "../../tools/desktop/template" .

  mkdir -p "./$PF_NAME$.app/Contents/Resources"
  mkdir -p "./$PF_NAME$.app/Contents/MacOS"

  cp -R "./template/macos/Info.plist" "./$PF_NAME$.app/Contents/Info.plist"
  cp -R "./template/macOS/icons.icns" "./$PF_NAME$.app/Contents/Resources/icons.icns"

  cp "$PF_KEY$.macos" "./$PF_NAME$.app/Contents/MacOS/$PF_KEY$"

  echo "signing desktop binary..."
  codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./$PF_NAME$.app/Contents/MacOS/$PF_KEY$"
  codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./$PF_NAME$.app"

  cp "./template/macos/appdmg.config.json" "./appdmg.config.json"

  echo "building macOS DMG..."
  appdmg "appdmg.config.json" "./$PF_KEY$_desktop_${TGT}_macos_x86_64.dmg"
  zip -r "$PF_KEY$_desktop_${TGT}_macos_x86_64.zip" "./$PF_NAME$.app"

  echo "building Linux zip..."
  zip "$PF_KEY$_desktop_${TGT}_linux_x86_64.zip" "./$PF_KEY$"

  echo "building Windows zip..."
  zip "$PF_KEY$_desktop_${TGT}_windows_x86_64.zip" "./$PF_KEY$.exe"

  mkdir -p "../../build/dist"
  mv "./$PF_KEY$_desktop_${TGT}_macos_x86_64.dmg" "../../build/dist"
  mv "./$PF_KEY$_desktop_${TGT}_macos_x86_64.zip" "../../build/dist"
  mv "./$PF_KEY$_desktop_${TGT}_linux_x86_64.zip" "../../build/dist"
  mv "./$PF_KEY$_desktop_${TGT}_windows_x86_64.zip" "../../build/dist"
fi
