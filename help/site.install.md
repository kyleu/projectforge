# Install Project Forge

You can install Project Forge in several ways:

## Go install

```
go install projectforge.dev/projectforge@latest
```

## Homebrew (macOS/Linux)

```
brew install kyleu/kyleu/projectforge
```

## Pre-built binaries

Download a release from GitHub and unzip the binary:

```
https://github.com/kyleu/projectforge/releases
```

## Docker

```
docker run -p 40000:40000 ghcr.io/kyleu/projectforge:latest
```

## First run

After installing:

```
projectforge create
./bin/dev.sh
```
