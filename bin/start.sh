#!/bin/bash

## Starts the app, without reloading on changes

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth
export projectforge_encryption_key=TEMP_SECRET_KEY

# include env file
if [ -f ".env" ]; then
	export $(cat .env | grep -v "#" | xargs)
fi

make clean
make build
./build/debug/projectforge -v
