#!/bin/bash
set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir

echo "running [build-machine]...";
tart run --dir=data:./data --dir={{{ .Exec }}}:../.. {{{ .Exec }}}-build-machine

