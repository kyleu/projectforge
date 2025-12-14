# Desktop

The **`desktop`** module enables your application to run as a native desktop application, across Windows, macOS, and Linux, using system webviews.

## Overview

This module transforms your web application into a cross-platform desktop application using:

- **Native Desktop Apps**: Creates platform-specific executables for Windows, macOS, and Linux
- **System Webview Integration**: Uses the operating system's built-in webview for rendering
- **Embedded Server**: Bundles your web server as a library within the desktop application
- **Cross-Platform Builds**: Automated build process supporting multiple architectures

## Key Features

### Multi-Platform Support

- **Windows**: Uses Edge WebView2 for modern web standards support
- **macOS**: Leverages WKWebView with native performance
- **Linux**: Utilizes WebKitGTK for consistent rendering

### Performance Benefits

- Native OS integration with system webviews
- No separate browser runtime required
- Minimal memory footprint compared to Electron
- Direct access to system resources

### Developer Experience

- Single codebase for web and desktop versions
- Automatic server lifecycle management
- Cross-compilation via Docker containers
- Consistent UI/UX across platforms

## Package Structure

### Build Infrastructure

- **`bin/`** - Desktop-specific build scripts and automation

  - Cross-platform compilation utilities
  - Package generation for different OS targets
  - Distribution and signing workflows

- **`tools/`** - Platform-specific build tools and dependencies
  - WebView runtime configurations
  - OS-specific compilation requirements
  - Packaging and installer generation

## Building Desktop Applications

### Prerequisites

- Docker (required for cross-compilation)

### Build Process

```bash
# Build desktop applications for all platforms
./bin/build/desktop.sh

```

### Cross-Compilation

Due to platform-specific dependencies and CGO requirements, desktop builds are performed in Docker containers with pre-configured toolchains for each target platform.

## Supported Platforms

| Platform | Architecture | WebView Technology | Notes                     |
| -------- | ------------ | ------------------ | ------------------------- |
| Windows  | x64, ARM64   | Edge WebView2      | Requires WebView2 runtime |
| macOS    | x64, ARM64   | WKWebView          | Native system component   |
| Linux    | x64, ARM64   | WebKitGTK          | Requires gtk3-devel       |

## Dependencies

The desktop module integrates with:

- **[WebView](https://github.com/webview/webview)** - Cross-platform webview library
- **Platform WebViews** - System-native rendering engines
- **Docker** - Cross-compilation environment
- **CGO** - Required for native webview bindings

## Limitations

- **CGO Dependency**: Cross-compilation requires Docker due to CGO requirements
- **WebView Availability**: Target systems must have compatible webview runtimes
- **Platform Differences**: Some web features may behave differently across webview implementations

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/desktop
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
