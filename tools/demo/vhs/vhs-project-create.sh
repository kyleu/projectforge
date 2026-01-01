#!/bin/bash

## Runs the project creation VHS script (https://github.com/charmbracelet/vhs)

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir"

if ! command -v vhs >/dev/null 2>&1; then
  echo "vhs is not installed. Please install it from https://github.com/charmbracelet/vhs"
  exit 1
fi

mkdir "$dir/tmp_project"
cd "$dir/tmp_project"

echo "Recording project creation..."
vhs < ../project-create.tape

cd "$dir"
rm -rf "$dir/tmp_project"
mkdir "$dir/tmp_project"
cd "$dir/tmp_project"

echo "Recording project creation (light mode)..."
vhs < ../project-create.light.tape

cd "$dir"
rm -rf "$dir/tmp_project"

echo "Completed recording"
