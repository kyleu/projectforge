# Project Forge WASM Server

This folder contains the templates used for running the WebAssembly build in a browser.

It uses the WebAssembly HTTP server as a ServiceWorker, allowing the application to be run offline.

To run it, execute `./bin/build/wasmserver.sh`, then host the files in `./build/wasm` using an HTTP server of your choosing.
