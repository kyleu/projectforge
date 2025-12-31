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
  echo "We're about to remove the homebrew [projectforge] installation, so [vhs] can reinstall it. Press enter to continue..."
  read -r
fi

brew remove -f kyleu/kyleu/projectforge

echo "Recording Homebrew installation..."
export HOMEBREW_NO_INSTALL_CLEANUP=1
vhs < homebrew-install.tape
echo "Completed recording"
