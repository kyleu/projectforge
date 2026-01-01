#!/bin/bash

## Runs the Golang installation VHS script (https://github.com/charmbracelet/vhs)

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

go_remove() {
  rm -f "${GOPATH}/bin/projectforge"
}
echo "Prepping Golang installation..."
go_remove
go install projectforge.dev/projectforge@latest
go_remove

echo "Recording Golang installation..."
vhs < golang-install.tape

echo "Recording Golang installation (light mode)..."
go_remove
vhs < golang-install.light.tape

echo "Completed recording"
