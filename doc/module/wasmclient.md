# WASM Client

The **`wasmclient`** module provides WebAssembly (WASM) client capabilities for [Project Forge](https://projectforge.dev) applications. It enables you to compile portions of your Go application to WebAssembly and run it in web browsers, providing a bridge between server-side Go code and client-side execution.

## Overview

This module extends Project Forge applications with:

- **WebAssembly Compilation**: Compile Go code to WASM for browser execution
- **JavaScript Bridge**: Seamless interaction between WASM and browser JavaScript
- **HTML Host Environment**: Ready-to-use HTML container for WASM applications
- **Build Tools**: Automated build scripts for WASM compilation and optimization

## Key Features

### WebAssembly Integration
- Go-to-WASM compilation with optimized build process
- Automatic JavaScript bridge generation
- DOM manipulation and browser API access
- File system and HTTP request capabilities

### Development Experience
- Hot-reload development server support
- Integrated with Project Forge build system
- Sandbox testing environment (when `sandbox` module is enabled)
- Comprehensive error handling and debugging

### Performance
- Optimized WASM binary size
- Efficient memory management
- Minimal JavaScript overhead
- Fast load times and execution

## Package Structure

### Core WASM Components

- **`app/wasm/main.go`** - WASM application entry point
  - Initialization and setup
  - Event loop management
  - Cleanup and resource management

- **`app/wasm/funcs.go`** - Exported WASM functions
  - JavaScript-callable Go functions
  - Type conversion helpers
  - Browser API integrations

### Build Infrastructure

- **`bin/build/wasmclient.sh`** - WASM build script
  - Go-to-WASM compilation
  - Binary optimization
  - Asset generation and packaging

### HTML Host (with Sandbox Module)

- **Testbed Environment** - Complete HTML host for testing
  - WASM loading and initialization
  - JavaScript bridge setup
  - Development debugging tools

## Usage

### Development Workflow

1. **Write Go code** in `app/wasm/main.go` and `app/wasm/funcs.go`
2. **Export functions** for JavaScript access using `//export` comments
3. **Build and test** using the provided build scripts
4. **Deploy** the generated WASM and JavaScript files

### Integration with Sandbox

When the `sandbox` module is enabled, a complete testbed environment is created that provides:
- HTML host page for your WASM application
- JavaScript debugging console
- Performance monitoring tools
- Interactive testing interface

## Configuration

The module works with standard Project Forge configuration and requires:
- Go 1.25+ with WebAssembly support
- Modern web browser with WASM support
- `sandbox` module (optional, for testbed environment)

## Build Process

The WASM build process:
1. Compiles Go source to WebAssembly binary
2. Generates JavaScript bridge code
3. Optimizes binary size and performance
4. Creates deployment-ready assets

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/wasmclient
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Sandbox Module](sandbox.md) - HTML testbed and debugging environment
- [Customization Guide](../customizing.md) - Advanced customization options
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
