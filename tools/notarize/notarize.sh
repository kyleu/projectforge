#!/bin/bash
# Content managed by Project Forge, see [projectforge.md] for details.

set -eo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

if [ "$PUBLISH_TEST" != "true" ]
then
  time gon "./tools/notarize/gon.amd64.hcl"
  time gon "./tools/notarize/gon.arm64.hcl"
  time gon "./tools/notarize/gon.all.hcl"
fi
