#!/bin/bash

## Runs all the tests

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

go test -v ./app/... | grep -v ^\?
