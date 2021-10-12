# Installation

## Pre-built binaries
Download any package from the [release page](https://{{{ .Package }}}/releases).

### Homebrew
```
brew install {{{ .Info.Org }}}/{{{ .Info.Org }}}/{{{ .Key }}} 
```

### deb, rpm and apk packages
Download the .deb, .rpm or .apk packages from the [release page](https://{{{ .Package }}}/releases) and install them with the appropriate tools.

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
git clone https://{{{ .Package }}}
cd {{{ .Key }}}
go mod tidy
make build
./build/debug/{{{ .Key }}} --help
```
