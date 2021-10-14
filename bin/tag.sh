#!/bin/bash

## Tags the git repo

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

TGT=$1
[ "$TGT" ] || (echo "must provide one argument, like \"0.0.1\"" && exit)

find . -type f -name "main.go" -print0 | xargs -0 sed -i '' -e "s/version = \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/version = \"${TGT}\"/g"
find . -type f -name "lib.go" -print0 | xargs -0 sed -i '' -e "s/Version: \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/Version: \"${TGT}\"/g"
find . -type f -name ".projectforge.json" -print0 | xargs -0 sed -i '' -e "s/\\\"version\\\": \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/\"version\": \"${TGT}\"/g"

cd tools/notarize
find . -type f -name "gon.*.hcl" -print0 | xargs -0 sed -i '' -e "s/_[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]_/_${TGT}_/g"
cd ../..

make build

git add .
git commit -m "v${TGT}" || true

git tag "v${TGT}"

git push
git push --tags
