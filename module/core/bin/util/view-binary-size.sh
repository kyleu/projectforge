#!/bin/bash

## Visualizes space usage for the release binary

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/../..

make build-release
go tool nm -size build/release/$PF_EXECUTABLE$ | c++filt > ./tmp/nm.txt

