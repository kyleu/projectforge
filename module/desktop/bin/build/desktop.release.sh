#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || TGT="v0.0.0"

if [ "$XSKIP_DESKTOP" != "true" ]
then
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
  mv dist/linux_amd64/{{{ .Key }}} "{{{ .Key }}}"
  mv dist/windows_amd64/{{{ .Key }}} "{{{ .Key }}}.exe"
  rm -rf "dist"

  # macOS
  cp -R "../../tools/desktop/template" .

  mkdir -p "./{{{ .Name }}}.app/Contents/Resources"
  mkdir -p "./{{{ .Name }}}.app/Contents/MacOS"

  cp -R "./template/macos/Info.plist" "./{{{ .Name }}}.app/Contents/Info.plist"
  cp -R "./template/macOS/icons.icns" "./{{{ .Name }}}.app/Contents/Resources/icons.icns"

  cp "{{{ .Key }}}.macos" "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"

  echo "signing desktop binary..."
  codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./{{{ .Name }}}.app/Contents/MacOS/{{{ .Key }}}"
  codesign  -f --options=runtime --verbose=4 --deep --force --strict -s 'Developer ID Application: Kyle Unverferth (C6S478FYLD)' "./{{{ .Name }}}.app"

  cp "./template/macos/appdmg.config.json" "./appdmg.config.json"

  echo "building macOS DMG..."
  appdmg "appdmg.config.json" "./{{{ .Key }}}_desktop_${TGT}_macos_x86_64.dmg"
  zip -r "{{{ .Key }}}_desktop_${TGT}_macos_x86_64.zip" "./{{{ .Name }}}.app"

  echo "building Linux zip..."
  zip "{{{ .Key }}}_desktop_${TGT}_linux_x86_64.zip" "./{{{ .Key }}}"

  echo "building Windows zip..."
  zip "{{{ .Key }}}_desktop_${TGT}_windows_x86_64.zip" "./{{{ .Key }}}.exe"

  mkdir -p "../../build/dist"
  mv "./{{{ .Key }}}_desktop_${TGT}_macos_x86_64.dmg" "../../build/dist"
  mv "./{{{ .Key }}}_desktop_${TGT}_macos_x86_64.zip" "../../build/dist"
  mv "./{{{ .Key }}}_desktop_${TGT}_linux_x86_64.zip" "../../build/dist"
  mv "./{{{ .Key }}}_desktop_${TGT}_windows_x86_64.zip" "../../build/dist"
fi
