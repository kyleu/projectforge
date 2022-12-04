# Installation

## Pre-built binaries
Download any package from the [release page]({{{ .Info.Sourcecode }}}/releases).
{{{ if .Build.Homebrew }}}
### Homebrew
```shell
brew install {{{ .Info.Org }}}/{{{ .Info.Org }}}/{{{ .Key }}}
```
{{{ end }}}{{{ if .Build.NFPMS }}}
### deb, rpm and apk packages
Download the .deb, .rpm or .apk packages from the [release page]({{{ .Info.Sourcecode }}}/releases) and install them with the appropriate tools.
{{{ end }}}
## Running with Docker
```shell
docker run -p {{{ .Port }}}:{{{ .Port }}} ghcr.io/{{{ .Info.Org }}}/{{{ .Key }}}:latest
docker run -p {{{ .Port }}}:{{{ .Port }}} ghcr.io/{{{ .Info.Org }}}/{{{ .Key }}}:latest-debug
```

## Built from source

### go install
```shell
go install {{{ .Package }}}@latest
```

### Source code

If you want to contribute to the project, please follow the steps on our [contributing guide](contributing).

If you just want to build from source for whatever reason, follow these steps:

```shell
git clone {{{ .Info.Sourcecode }}}
cd {{{ .Key }}}
make build
./build/debug/{{{ .Key }}} --help
```

A script has been provided at `./bin/dev.sh` that will auto-reload when the project changes.

Note that you may need to run `./bin/bootstrap.sh` before building the project for the first time.
