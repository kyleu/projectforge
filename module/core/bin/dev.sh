#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth

ulimit -n 2048
air
