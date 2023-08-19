<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# WebAssembly

This is a module for [Project Forge](https://projectforge.dev). It allows you to build your http server, but load it as a WebAssembly module. 

https://github.com/kyleu/projectforge/tree/master/module/wasmserver

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

To use your app with only WebAssembly, no server process required, run `./bin/build/wasmserver.sh`, then host the files in ./tools/wasmserver in an HTTP server (sadly, it won't work with `file://` urls).
