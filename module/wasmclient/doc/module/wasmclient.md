# WASM Client

This is a module for [Project Forge](https://projectforge.dev). Provides a WebAssembly library and HTML host for a custom WASM application

https://github.com/kyleu/projectforge/tree/main/module/wasm

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

A WASM library is available in `./app/wasm`. If the `sandbox` module is enabled, a testbed will be created that provides an HTML host for the application.

A main method is generated in `/app/wasm/main.go`, and exported functions are defined in `/app/wasm/funcs.go`.
To build your wasm project, run `/bin/build/wasmclient.sh`.
