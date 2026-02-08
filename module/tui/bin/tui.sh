#!/bin/bash

## Builds and runs the TUI, then re-runs on file changes.
##
## Usage:
##   ./bin/tui.sh
##
## Requires:
##   - make
##   - watchexec in PATH
##
## Notes:
##   - Watches the repository root
##   - Ignores build artifacts to avoid watch loops

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd make "install Xcode command line tools or GNU make"
require_cmd watchexec "install from https://github.com/watchexec/watchexec"

echo "Watching for changes and running TUI..."
watchexec \
  --on-busy-update=restart \
  --project-origin . \
  --no-discover-ignore \
  --ignore-file .gitignore \
  --exts go,html,md,mod{{{ if .HasModule "database" }}},sql{{{ end }}},js,css{{{ if .HasModule "graphql" }}},graphql,schema{{{ end }}} \
  --watch . \
  --ignore .git \
  --ignore "build" \
  --ignore "build/**" \
  --ignore "client/**" \
  --ignore "gen/**" \
  --ignore "tools/**" \{{{ if .HasModule "notebook" }}}
  --ignore "notebook/**" \{{{ end }}}{{{ if .HasModule "playwright" }}}
  --ignore "test/playwright/**" \{{{ end }}}
  --ignore "data/**" \
  --ignore "module/**" \
  --ignore "testproject/**" \
  --ignore "tmp/**" \
  --ignore "assets/module/**" \
  --ignore "**/*.html.go" \{{{ if .HasModule "database" }}}
  --ignore "**/*.sql.go" \{{{ end}}}
  --ignore "**/*_test.go" \
  --shell=bash \
  -- "make build && ./build/debug/projectforge tui"
