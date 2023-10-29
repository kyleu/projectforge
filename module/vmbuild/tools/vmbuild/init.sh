#!/bin/bash
set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $dir

echo "initializing tart virtual machine...";
tart clone ghcr.io/cirruslabs/macos-sonoma-base:latest {{{ .Exec }}}-build-machine

tart set {{{ .Exec }}}-build-machine --cpu 6 --memory 4096
