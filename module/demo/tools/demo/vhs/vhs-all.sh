#!/bin/bash

## Runs the VHS scripts (https://github.com/charmbracelet/vhs)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"

if ! command -v vhs >/dev/null 2>&1; then
  echo "vhs is not installed. Please install it from https://github.com/charmbracelet/vhs"
  exit 1
fi

./vhs-golang-install.sh "$@"{{{ if .Build.Homebrew }}}
./vhs-homebrew-install.sh "$@"{{{ end }}}
