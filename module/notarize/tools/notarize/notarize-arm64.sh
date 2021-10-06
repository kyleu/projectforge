#!/bin/bash

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

if [ "$PUBLISH_TEST" != "true" ]
then
  time gon ./tools/notarize/gon.arm64.hcl
fi
