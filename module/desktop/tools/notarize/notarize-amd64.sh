#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

TGT=$1
[ "$TGT" ] || (echo "must provide one argument, like \"0.0.1\"" && exit)

sed -i '' "s/[v]*[01]\.[0-9]*[0-9]\.[0-9]*[0-9][-SNAPSHOT]*/$TGT/g" ./tools/notarize/gon.amd64.hcl
time gon ./tools/notarize/gon.amd64.hcl
