# Model Context Protocol (MCP)

The **`mcp`** module provides a complete [Model Context Protocol](https://modelcontextprotocol.io) server implementation for your application. MCP is a standardized protocol that enables Large Language Models to securely connect to external tools and data sources.

## Overview

This module provides a complete MCP protocol implementation with tools, resources, and prompts. A UI is provided for testing, and the server can be exposed over HTTP and called via command line

## Integration

**`lib/mcpserver/`** - MCP server implementation. Add your tools and resources here.

### Tool Development

Tools are defined using the `Tool` struct:

```go
type ToolHandler func(ctx context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (string, error)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Icon        string          `json:"icon,omitempty"`
	Args        util.FieldDescs `json:"args,omitempty"`
	Fn          ToolHandler     `json:"-"`
}
```

## Usage

### CLI Integration

Configure your MCP client to use your application as an MCP server via CLI:

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

...or via HTTP:

```json
{
  "mcpServers": {
    "myapp": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "http://localhost:port/mcp"]
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
