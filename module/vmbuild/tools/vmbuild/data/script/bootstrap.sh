#!/bin/bash
set -euo pipefail
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$dir/.."

echo "### Initializing macOS build machine...";

echo '### Setting maxfiles...'
cp assets/limit.maxfiles.plist ~/limit.maxfiles.plist
sudo mv ~/limit.maxfiles.plist /Library/LaunchDaemons/limit.maxfiles.plist
sudo chown root:wheel /Library/LaunchDaemons/limit.maxfiles.plist
sudo chmod 0644 /Library/LaunchDaemons/limit.maxfiles.plist

echo '### Disabling spotlight...'
sudo mdutil -a -i off

echo "### Installing Safari driver..."
sudo safaridriver --enable

echo "### Installing Rosetta..."
sudo softwareupdate --install-rosetta --agree-to-license

echo "### Installing Xcode..."
echo 'export PATH=/usr/local/bin/:$PATH' >> ~/.zprofile
source ~/.zprofile
wget --quiet https://github.com/RobotsAndPencils/xcodes/releases/latest/download/xcodes.zip
unzip xcodes.zip
rm xcodes.zip
chmod +x xcodes
sudo mkdir -p /usr/local/bin/
sudo mv xcodes /usr/local/bin/xcodes
xcodes version
wget --quiet https://storage.googleapis.com/xcodes-cache/Xcode_15.xip
xcodes install 15 --experimental-unxip --path $PWD/Xcode_15.xip
sudo rm -rf ~/.Trash/*
xcodes select 15
xcodebuild -downloadAllPlatforms
xcodebuild -runFirstLaunch

echo '### Installing Homebrew...'
NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
echo "export LANG=en_US.UTF-8" >> ~/.zprofile
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
echo "export HOMEBREW_NO_AUTO_UPDATE=1" >> ~/.zprofile
echo "export HOMEBREW_NO_INSTALL_CLEANUP=1" >> ~/.zprofile
source ~/.zprofile
brew --version
brew update
brew install cmake curl gcc wget

echo "### Installing Node.js..."
brew install node
node --version
npm install --global npm@latest
npm install --global yarn
yarn --version

echo "### Installing Golang..."
brew install go
go version
echo "export PATH=$PATH:${HOME}/go/bin" >> ~/.zprofile
source ~/.zprofile

echo "### Installing Go dependencies..."
go install github.com/cosmtrek/air@latest
go install github.com/valyala/quicktemplate/qtc@latest
go install gotest.tools/gotestsum@latest
go install mvdan.cc/gofumpt@latest
go install github.com/goreleaser/goreleaser@latest

echo "### Installing [~/projects] link..."
ln -s "/Volumes/My Shared Files" ~/projects

echo "### All Done!"
