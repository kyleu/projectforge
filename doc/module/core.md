# Core

The **`core`** module is the foundational module for [Project Forge](https://projectforge.dev) applications. It provides essential utilities, infrastructure components, and common patterns that form the backbone of every application managed by Project Forge.

## Overview

This module is **required** for all Project Forge applications and provides:

- **CLI Infrastructure**: Command-line interface framework and built-in commands
- **HTTP Server**: High-performance web server with routing, middleware, and content negotiation
- **Observability**: Comprehensive logging, metrics, and distributed tracing
- **UI Framework**: Themeable UI components, navigation, and filtering
- **Utilities**: Performance-optimized helper functions and common patterns

## Key Features

### Performance
- Sub-second build times
- <20KB total payload
- <1ms response times
- Zero-allocation utilities where possible

### Developer Experience
- Live reload development server
- Comprehensive CLI tooling
- Built-in profiling and debugging
- Extensive configuration options

### Observability
- OpenTelemetry distributed tracing
- Prometheus metrics collection
- Structured logging with multiple formatters
- Built-in health checks and diagnostics

### UI Framework
- Multiple built-in themes (light/dark modes)
- Responsive navigation components
- Advanced filtering and sorting
- Progressive enhancement (works without JavaScript)

## Package Structure

### Core Infrastructure

- **`cmd/`** - CLI command framework and built-in commands
  - Application lifecycle management
  - Development server commands
  - Build and deployment utilities

- **`controller/`** - HTTP request handlers and middleware
  - Authentication and authorization
  - Content negotiation (JSON, HTML, CSV, XML, YAML)
  - Error handling and recovery
  - CORS and security headers

### Libraries

- **`lib/filter/`** - Advanced filtering and sorting for UI and APIs
  - Multi-column sorting
  - Type-aware filtering (strings, numbers, dates, booleans)
  - URL-based filter persistence

- **`lib/log/`** - Structured logging infrastructure
  - Multiple output formats (JSON, console, file)
  - Custom zap loggers and appenders
  - Request correlation and tracing integration

- **`lib/menu/`** - Navigation and breadcrumb systems
  - Dynamic menu generation
  - Role-based menu filtering
  - Breadcrumb trail management

- **`lib/telemetry/`** - Observability and monitoring
  - OpenTelemetry trace collection
  - Prometheus metrics integration
  - Custom metrics and spans
  - Performance profiling utilities

- **`lib/theme/`** - UI theming and customization
  - Multiple built-in themes
  - Dark/light mode support
  - Custom CSS property integration
  - Dynamic theme switching

- **`lib/user/`** - User management and permissions
  - User account management
  - Role-based access control
  - Session management
  - Permission checking utilities

### Utilities

- **`util/`** - Performance-optimized helper functions
  - String manipulation and formatting
  - Data structure utilities
  - File and path operations
  - HTTP client helpers
  - Environment variable handling

## Configuration

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/core
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Customization Guide](../customizing.md) - Advanced customization options
- [Configuation Variables](../running.md) - Available environment variables for configuration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
