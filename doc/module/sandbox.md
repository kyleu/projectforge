# Sandbox

The **`sandbox`** module provides interactive testing environments for [Project Forge](https://projectforge.dev) applications. It enables developers to create custom playgrounds for experimenting with application features, testing new functionality, and prototyping code in a safe, isolated environment.

## Overview

This module provides:

- **Interactive Testbeds**: Customizable sandbox environments for testing code
- **Dynamic Arguments**: Form-based parameter collection for sandbox functions
- **Admin Integration**: Seamless integration with the admin interface
- **Extensible Framework**: Easy addition of new sandbox environments
- **Safe Execution**: Isolated execution context for experimental code

## Key Features

### Development & Testing
- Interactive web interface for running test functions
- Dynamic form generation for function parameters
- Real-time execution and result display
- Integration with application telemetry and logging

### Customization
- Easy creation of new sandbox environments
- Flexible argument specification using field descriptors
- Custom result rendering and visualization
- Conditional module integration (e.g., WASM support)

### Safety & Monitoring
- Isolated execution context for experimental code
- Comprehensive logging and tracing integration
- Error handling and recovery
- Performance monitoring with telemetry spans

## Package Structure

### Core Framework

- **`lib/sandbox/`** - Core sandbox infrastructure
  - **`sandbox.go`** - Main sandbox framework and types
  - **`testbed.go`** - Default testbed implementation
  - Dynamic sandbox registration and discovery
  - Menu integration for admin interface

### Web Interface

- **`controller/clib/`** - HTTP handlers for sandbox operations
  - **`sandbox.go`** - Web interface controllers
  - Dynamic parameter collection
  - Result rendering and display
  - Error handling and user feedback

### Templates

- **`views/vsandbox/`** - HTML templates for sandbox UI
  - Sandbox listing and navigation
  - Parameter input forms
  - Result display and visualization

## Usage

### Creating a New Sandbox

Add a new sandbox by creating a `Sandbox` struct, and registering it in the `AllSandboxes` slice:

```go
var mySandbox = &Sandbox{
    Key:   "my-test",
    Title: "My Test Environment", 
    Icon:  "beaker",
    Args: util.FieldDescs{
        {Key: "message", Title: "Message", Type: "string", Default: "Hello"},
        {Key: "count", Title: "Count", Type: "int", Default: "5"},
    },
    Run: onMyTest,
}

func onMyTest(ctx context.Context, st *app.State, args util.ValueMap, logger util.Logger) (any, error) {
    message := args.GetStringOr("message", "default")
    count := args.GetIntOr("count", 1)
    
    result := util.ValueMap{
        "message": message,
        "repeated": strings.Repeat(message+" ", count),
        "timestamp": time.Now(),
    }
    
    return result, nil
}
```

### Accessing Sandboxes

Sandboxes are accessible through the admin web interface:

1. Navigate to `/admin/sandbox` to see all available sandboxes
2. Click on a specific sandbox to access its interface
3. Fill in required parameters (if any)
4. Execute the sandbox function and view results

### URL Structure

- **List**: `/admin/sandbox` - Shows all available sandboxes
- **Run**: `/admin/sandbox/{key}` - Executes a specific sandbox
- **Parameters**: Query parameters or form submission for sandbox arguments

## Integration

### Menu System
The sandbox module automatically registers with the admin menu:
- Main menu item: "Sandboxes" 
- Submenu items: Individual sandboxes
- Dynamic menu generation based on available sandboxes

### Telemetry
All sandbox executions are traced with OpenTelemetry:
- Span name: `sandbox.{key}`
- Automatic performance monitoring
- Error tracking and logging

### Error Handling
- Graceful error handling for missing sandboxes
- Parameter validation and missing argument detection
- User-friendly error messages in the web interface

## Advanced Features

### Custom Result Rendering
Specific sandboxes can have custom result templates:
- Default: Generic result display
- Custom: Specialized rendering (e.g., `testbed` template)
- Conditional: Module-dependent rendering (e.g., WASM support)

### Parameter Types
Supports various parameter types through field descriptors:
- **string**: Text input
- **int**: Numeric input  
- **bool**: Checkbox input
- **select**: Dropdown selection
- **textarea**: Multi-line text input

### Security Considerations
- Admin-only access through existing authentication
- Isolated execution context
- Input validation and sanitization
- Safe error handling without exposing internals

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/sandbox
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Core Module Documentation](core.md) - Foundation framework
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
- [Customizing Guide](../customizing.md) - Advanced customization options
