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

BREW_PATH="kyleu/kyleu/projectforge"

brew_remove() {
  brew remove -f $BREW_PATH
}

echo "Prepping Homebrew installation..."
export HOMEBREW_NO_INSTALL_CLEANUP=1
brew_remove
brew update
brew install $BREW_PATH
brew_remove

echo "Recording Homebrew installation..."
vhs < homebrew-install.tape

echo "Recording Homebrew installation (light mode)..."
brew_remove
vhs < homebrew-install.light.tape

echo "Completed recording"
