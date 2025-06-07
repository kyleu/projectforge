# Model Context Protocol (MCP)

The **`mcp`** module provides a complete [Model Context Protocol](https://modelcontextprotocol.io) server implementation for [Project Forge](https://projectforge.dev) applications. MCP is a standardized protocol that enables Large Language Models to securely connect to external tools and data sources.

## Overview

This module provides:

- **MCP Server**: Complete MCP protocol implementation with tool registration
- **CLI Interface**: Command-line interface for serving MCP requests via stdio
- **Web Admin Interface**: Browser-based tool testing and management interface
- **Tool Framework**: Extensible framework for creating custom MCP tools
- **Example Tools**: Built-in example tools to demonstrate usage patterns

## Key Features

### Protocol Compliance

- Full MCP specification compliance
- Secure request/response handling
- Error handling and validation
- Tool discovery and registration

### Developer Experience

- Simple tool registration API
- Built-in example tools as templates
- Web-based testing interface
- Comprehensive logging and debugging

### Integration

- CLI command: `<app> mcp` for MCP client integration
- HTTP admin interface at `/admin/mcp`
- Seamless integration with Project Forge applications
- Support for custom tool development

## Package Structure

### Core Components

- **`cmd/mcp.go`** - CLI command implementation

  - MCP server initialization
  - Stdio protocol handling
  - Context management

- **`lib/mcpserver/`** - MCP server implementation

  - `server.go` - Core MCP server and tool registration
  - `tool.go` - Tool definition and handler framework
  - `example.go` - Example tools for demonstration

- **`controller/clib/mcp.go`** - Web admin interface
  - Tool listing and management
  - Interactive tool testing
  - Request/response inspection

### Tool Development

Tools are defined using the `Tool` struct:

```go
type Tool struct {
    Name        string
    Description string
    InputSchema map[string]any
    Handler     func(context.Context, map[string]any) (any, error)
}
```

## Usage

### CLI Integration

Configure your MCP client to use your application as an MCP server:

```json
{
  "mcpServers": {
    "myapp": {
      "command": "/path/to/myapp",
      "args": ["mcp"]
    }
  }
}
```

### Adding Custom Tools

Register custom tools in your application:

```go
// In your application initialization
server, err := mcpserver.NewServer(ctx, version)
if err != nil {
    return err
}

customTool := &mcpserver.Tool{
    Name:        "custom-tool",
    Description: "Description of what this tool does",
    InputSchema: map[string]any{
        "type": "object",
        "properties": map[string]any{
            "input": map[string]any{
                "type": "string",
                "description": "Input parameter description",
            },
        },
        "required": []string{"input"},
    },
    Handler: func(ctx context.Context, params map[string]any) (any, error) {
        // Tool implementation
        return result, nil
    },
}

server.AddTools(customTool)
```

### Web Admin Interface

Access the web admin interface at `/admin/mcp` to:

- View registered tools
- Test tool execution
- Inspect request/response data
- Debug tool implementations

## Example Tools

The module includes example tools to demonstrate usage:

- **Random Number Generator** - Generates random numbers within specified ranges
- **Generated Tools** - Tools created from Project Forge export configurations (if enabled)

## Configuration

The MCP module uses standard Project Forge configuration patterns:

- **Debug Mode**: Enable detailed logging with `DEBUG=true`
- **Build Info**: Automatic version and build information integration
- **Tool Registration**: Tools are registered during server initialization

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/mcp
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Model Context Protocol Specification](https://modelcontextprotocol.io) - Official MCP documentation
- [MCP Go Library](https://github.com/mark3labs/mcp-go) - Underlying Go implementation
- [Project Forge Documentation](https://projectforge.dev) - Complete Project Forge documentation
- [Customization Guide](../customizing.md) - Advanced customization options
