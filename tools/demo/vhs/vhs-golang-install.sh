#!/bin/bash

## Runs the VHS scripts (https://github.com/charmbracelet/vhs)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"

if ! command -v vhs >/dev/null 2>&1; then
  echo "vhs is not installed. Please install it from https://github.com/charmbracelet/vhs"
  exit 1
fi

prompt=true
for arg in "$@"; do
  if [[ "$arg" == "-f" ]]; then
    prompt=false
    break
  fi
done
if "$prompt"; then
  echo "We're about to remove the golang-installed [projectforge] binary, so [vhs] can reinstall it. Press enter to continue..."
  read -r
fi

rm -f "${GOPATH}/bin/projectforge"

echo "Recording Golang installation..."
vhs < golang-install.tape
echo "Completed recording"
