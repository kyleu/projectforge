#!/bin/bash

## Starts the app, reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

# $PF_SECTION_START(keys)$
# $PF_SECTION_END(keys)$

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth
export {{{ .Key }}}_encryption_key=TEMP_SECRET_KEY

ulimit -n 2048
air
