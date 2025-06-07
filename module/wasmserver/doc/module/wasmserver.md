# WASM Server

The **`wasmserver`** module for [Project Forge](https://projectforge.dev) enables building your HTTP server application as a WebAssembly module that runs in web browsers using Service Workers. This allows your Go application to run entirely client-side without requiring a traditional server infrastructure.

## Overview

This module transforms your Project Forge application into a WebAssembly binary that can be executed in web browsers, enabling:

- **Client-side execution**: Full application logic runs in the browser
- **Service Worker integration**: Intercepts HTTP requests within the browser
- **Serverless deployment**: Host static files without server infrastructure
- **Offline capabilities**: Applications can work without network connectivity

## Key Features

### WebAssembly Compilation
- Compiles Go HTTP server to WASM binary
- Maintains full application functionality
- Uses Go's built-in WASM target (`GOOS=js GOARCH=wasm`)

### Service Worker Integration
- Automatically generates Service Worker wrapper
- Intercepts HTTP requests within browser context
- Provides seamless request/response handling
- Enables progressive web app capabilities

### Static File Generation
- Creates deployable static file structure
- Includes necessary WASM runtime files
- Generates HTML wrapper for application bootstrap
- Provides hosting-ready file structure

## Build Process

### Building WASM Server

```bash
# Build the WASM server
./bin/build/wasmserver.sh

# Files are generated in ./tools/wasmserver/
```

### Generated Files

The build process creates:

- **`main.wasm`** - Compiled Go application
- **`wasm_exec.js`** - Go WASM runtime support
- **`sw.js`** - Service Worker implementation
- **`index.html`** - Application bootstrap page
- **Static assets** - CSS, JavaScript, and other resources

## Deployment

### Requirements

- **HTTPS required**: Service Workers require secure context (except `localhost`)
- **Web server**: Cannot use `file://` protocol, requires HTTP server
- **Modern browser**: Supports WebAssembly and Service Workers

### Hosting Options

```bash
# Local development
cd tools/wasmserver
python3 -m http.server 8080

# Or using Node.js
npx serve .

# Or any static file server
```

### Production Deployment

Deploy the `tools/wasmserver` directory contents to any static hosting service:
- Netlify, Vercel, GitHub Pages
- AWS S3 + CloudFront
- Traditional web hosting
- CDN services

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
