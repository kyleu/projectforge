#!/bin/bash

## Downloads and installs the Go libraries and tools needed in other scripts

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir/..

go get -u github.com/cosmtrek/air
go get -u github.com/pyros2097/go-embed
go get -u github.com/valyala/quicktemplate
go get -u github.com/valyala/quicktemplate/qtc
go mod download