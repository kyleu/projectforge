# Contributing

Source code is available at [GitHub]({{{ .Info.Sourcecode }}}).

## Setup

```
git clone https://github.com/kyleu/projectforge
cd projectforge
./bin/bootstrap.sh
```

## Build and run

```
make build
./build/debug/projectforge --help
./bin/dev.sh
```

## Before submitting changes

```
./bin/format.sh
./bin/check.sh
./bin/test.sh
make build
```

## Notes

- Go and Node are required for full builds.
- The dev server auto-rebuilds templates and client assets.
