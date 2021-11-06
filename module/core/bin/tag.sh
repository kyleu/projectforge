#!/bin/bash

## Tags the git repo using the first argument or the incremented minor version

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

TGT=${1-none}
if [[ $TGT == "none" ]]; then
  TGT=$(git describe --tags | sed -e 's/v//g')
  TGT=$(echo ${TGT} | awk -F. -v OFS=. '{$NF++;print}')
fi
if [[ ${TGT:0:1} == "v" ]]; then
  TGT = "${TGT:1}"
fi

echo $TGT

find . -type f -name "main.go" -print0 | xargs -0 sed -i '' -e "s/version = \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/version = \"${TGT}\"/g"{{{ if .UsesLib }}}
find . -type f -name "lib.go" -print0 | xargs -0 sed -i '' -e "s/Version: \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/Version: \"${TGT}\"/g"{{{ end }}}{{{ if .Build.Notarize }}}
find . -type f -name "gon.*.hcl" -print0 | xargs -0 sed -i '' -e "s/\\_[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]_/_${TGT}\\_/g"{{{ end }}}
find . -type f -name ".projectforge.json" -print0 | xargs -0 sed -i '' -e "s/\\\"version\\\": \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/\"version\": \"${TGT}\"/g"

make build

git add .
git commit -m "v${TGT}" || true

git tag "v${TGT}"

git push
git push --tags
