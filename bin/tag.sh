#!/bin/bash

## Bumps version, builds, commits, tags, and pushes.
##
## Usage:
##   ./bin/tag.sh [version|tag] [--dry-run] [-y|--yes]
##
## Arguments:
##   version|tag  If omitted, increments the patch segment of the latest vX.Y.Z tag.
##   --dry-run    Print intended actions and exit without changing the repo.
##   -y, --yes    Skip the confirmation prompt.
##
## Requires:
##   - Clean git working tree
##   - git, sed, awk, make
##
## Notes:
##   - Updates versions in source code
##   - Creates a commit, creates a tag, and pushes commits and tags

set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command '$1' not found${2:+ ($2)}" >&2
    exit 1
  fi
}

require_cmd git "install Git from https://git-scm.com/downloads"
require_cmd sed "install GNU sed (or use macOS default)"
require_cmd awk "install awk (typically preinstalled)"
require_cmd make "install Xcode Command Line Tools or build-essential"

auto_yes=false
dry_run=false
TGT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -y|--yes) auto_yes=true; shift;;
    --dry-run) dry_run=true; shift;;
    --) shift; break;;
    -*)
      echo "unknown option: $1" >&2
      exit 1
      ;;
    *)
      if [[ -z "$TGT" ]]; then
        TGT="$1"
        shift
      else
        echo "unexpected argument: $1" >&2
        exit 1
      fi
      ;;
  esac
done

# uncommitted file checks
git update-index -q --ignore-submodules --refresh
err=0

if ! git diff-files --quiet --ignore-submodules --
then
    echo >&2 "you have unstaged changes."
    git diff-files --name-status -r --ignore-submodules -- >&2
    err=1
fi

if ! git diff-index --cached --quiet HEAD --ignore-submodules --
then
    echo >&2 "your index contains uncommitted changes."
    git diff-index --cached --name-status -r --ignore-submodules HEAD -- >&2
    err=1
fi

if [ $err = 1 ]
then
    echo >&2 "please commit or stash them."
    exit 1
fi
# end checks

git fetch --tags

TGT=${TGT:-none}
if [[ $TGT == "none" ]]; then
  TGT=$(git describe --match "v[0-9]*" --tags --abbrev=0 | sed -e 's/v//g')
  TGT=$(echo "${TGT}" | awk -F. -v OFS=. '{$NF++;print}')
fi
if [[ ${TGT:0:1} == "v" ]]; then
  TGT="${TGT:1}"
fi

echo "$TGT"

if $dry_run; then
  echo "dry run: would update version strings (if numeric), build, commit, tag, and push"
  exit 0
fi

if ! $auto_yes; then
  read -r -p "Tag ${TGT} and push to origin? [y/N] " confirm
  case "$confirm" in
    [yY][eE][sS]|[yY]) ;;
    *) echo "aborted"; exit 0;;
  esac
fi

pat="^[0-9]"
if [[ $TGT =~ $pat ]]; then
  sed -i.bak -e "s/version = \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/version = \"${TGT}\"/g" ./main.go
  rm -f "./main.go.bak"
  sed -i.bak -e "s/Version: \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/Version: \"${TGT}\"/g" ./app/cmd/lib.go
  rm -f ./app/cmd/lib.go.bak
  sed -i.bak -e "s/\\_[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]_/_${TGT}\\_/g" ./tools/notarize/notarize.sh
  rm -f "./tools/notarize/notarize.sh.bak"
  sed -i.bak -e "s/\\\"version\\\": \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/\"version\": \"${TGT}\"/g" ./.projectforge/project.json
  rm -f "./.projectforge/project.json.bak"
  sed -i.bak -e "s/\\\"version\\\": \\\"[v]*[0-9]*[0-9]\.[0-9]*[0-9]\.[0-9]*[0-9]\\\"/\"version\": \"${TGT}\"/g" ./client/package.json
  rm -f "./client/package.json.bak"
fi

make build

git add .

if [[ $TGT =~ $pat ]]; then
  git commit -m "v${TGT}" || true
  git tag "v${TGT}"
else
  git commit -m "${TGT}" || true
  git tag "${TGT}"
fi

git push
git push --tags
