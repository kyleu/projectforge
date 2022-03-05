#!/bin/bash
set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir

cd ..
make build
cd ~/fevo-tech/goalkeeper
~/kyleu/projectforge/build/debug/projectforge debug
