#!/bin/bash

## Builds pfdb

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

if [ -f ".env" ]; then
	export "$(cat .env | grep -v "#" | xargs)"
fi

go build -gcflags "all=-N -l" -o build/pfdb .
