# WASM Server

The **`wasmserver`** module builds your existing HTTP server as a WebAssembly binary and runs it in the browser via a Service Worker. The result is a static bundle that behaves like your server routes, without a server process.

## Overview

This module provides:

- **Client-side execution**: Full server router runs in the browser
- **Service Worker integration**: Intercepts same-origin requests and dispatches them to Go
- **Static hosting**: Deploy to any static file host
- **Offline support**: Service Worker caching and offline behavior

## How It Works

1. `./bin/build/wasmserver.sh` compiles your server with `GOOS=js GOARCH=wasm`
2. Static wrapper files (`index.html`, `server.js`, `sw.js`, `wasm_exec.js`) are copied into `./build/wasm`
3. `server.js` registers `sw.js`, which loads `projectforge.wasm`
4. The Service Worker calls `goFetch` (from the Go WASM runtime) to handle requests with your router
5. The Service Worker checks for updated WASM binaries by ETag on install/activate and every 15 minutes

## Usage

```bash
# Build the WASM server bundle
./bin/build/wasmserver.sh

# Serve the output (uses npx http-server on port 40009)
./bin/wasmdev.sh
```

Then visit `http://localhost:40009` in a browser. After the Service Worker takes control (refresh once), normal routes are handled by your Go server code in the browser.

## Configuration

Configuration lives in `.projectforge/project.json`:

- `build.wasm`: Enables WASM outputs in release builds
- `info.homepage`: Used by `tools/wasmserver/sw.js` to skip intercepting the bootstrap URL; set it to your deployed origin, e.g. `https://example.com`

Template files you can customize:

- `tools/wasmserver/index.html`: bootstrap HTML
- `tools/wasmserver/server.js`: Service Worker registration
- `tools/wasmserver/sw.js`: request interception and WASM lifecycle

## CLI And URLs

Commands:

- `./bin/build/wasmserver.sh`: Build `./build/wasm`
- `./bin/build/wasmrelease.sh <version>`: Create `./build/dist/projectforge_<version>_wasm_html.zip`
- `./bin/wasmdev.sh`: Serve `./build/wasm` on `http://localhost:40009`
- `projectforge wasm`: Internal entrypoint used by the Service Worker (not for direct CLI use)

Static URLs served from `./build/wasm`:

- `/index.html`
- `/sw.js`
- `/server.js`
- `/projectforge.wasm`
- `/wasm_exec.js`
- `/assets/*`, `/favicon.ico`, `/logo.svg`

## Dependencies

- Go 1.25+ with WASM support
- A static web server (Node `npx http-server`, `python3 -m http.server`, or equivalent)
- `zip` if you run `./bin/build/wasmrelease.sh`
- A modern browser with WebAssembly + Service Worker support (HTTPS required in production)

## Examples

Local dev:

```bash
./bin/build/wasmserver.sh
./bin/wasmdev.sh
# open http://localhost:40009
```

Release bundle:

```bash
./bin/build/wasmrelease.sh 1.2.3
# upload build/dist/projectforge_1.2.3_wasm_html.zip to static hosting
```

## Limitations

### Browser Restrictions

- **No cookie setting**: Service Workers cannot set HTTP cookies
- **Limited headers**: Cannot set many security-sensitive headers
- **Cross-origin restrictions**: Subject to CORS policies
- **Storage limitations**: Relies on browser storage APIs

### Performance Considerations

- **Initial load time**: WASM binary download and compilation
- **Memory usage**: Runs within browser memory constraints
- **CPU limitations**: JavaScript execution context performance
- **File size**: WASM binaries can be larger than traditional JavaScript

### Functionality Gaps

- **Database connections**: No direct database connectivity
- **File system access**: Limited file system operations
- **Network operations**: Restricted to browser networking APIs
- **Process management**: No access to system processes

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/wasmserver
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [WebAssembly Documentation](https://webassembly.org/docs/) - WebAssembly specifications
- [Service Workers](https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API) - MDN Service Worker documentation
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
