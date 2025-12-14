# Model Context Protocol (MCP)

The **`mcp`** module provides a complete [Model Context Protocol](https://modelcontextprotocol.io) server implementation for your application. MCP is a standardized protocol that enables Large Language Models to securely connect to external tools and data sources.

## Overview

This module provides a complete MCP protocol implementation with tools, resources, and prompts. A UI is provided for testing, and the server can be exposed over HTTP and called via command line

## Usage

### Server Registration

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

If your tools support "http" server types, you can use:

```json
{
  "mcpServers": {
    "myapp": {
      "type": "http",
      "url": "http://localhost:port/mcp"
    }
  }
}
```

Otherwise, you can use the `mcp-remote` tool:

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

## Integration

**`lib/mcpserver/`** - MCP server implementation. Register your tools and resources here.

### Tools

```go
var CustomTool = &Tool{
	Name:        "custom_tool",
	Description: "Run a custom tool from " + util.AppName,
	Args: util.FieldDescs{
		{Key: "id", Description: "The id to us"},
	},
	Fn: func(_ context.Context, as *app.State, req mcp.CallToolRequest, args util.ValueMap, logger util.Logger) (any, error) {
		return "TODO", nil
	},
}

mcpServer, _ := mcpserver.NewServer(ctx, as, logger)
err := mcpServer.AddTools(as, logger, AllTools...)
```

### Resources (static)

```go
var CustomResource = mcpserver.NewTextResource("resource", "resource description", "resource://foo", "text/plain", "Hello!")

mcpServer, _ := mcpserver.NewServer(ctx, as, logger)
err := mcpServer.AddResources(as, logger, CustomResource)
```

### Resources (dynamic)

```go

var CustomResource = &ResourceTemplate{
	Name:        "custom_resource",
	Description: "Load a custom resource from " + util.AppName,
	URI:         "custom_resource://{id}/{path*}",
	Args: util.FieldDescs{
		{Key: "id", Description: "Some ID"},
		{Key: "path", Description: "A path to some file"},
	},
	Fn: func(_ context.Context, as *app.State, req mcp.ReadResourceRequest, args util.ValueMap, logger util.Logger) (string, string, any, error) {
		id := args.GetStringOpt("id")
		path := args.GetStringOpt("path")
		u := fmt.Sprintf("custom_resource://%s/%s", id, path)
		return u, "text/plain", []byte("Hello!"), nil
	},
}

mcpServer, _ := mcpserver.NewServer(ctx, as, logger)
err := mcpServer.AddResourceTemplates(as, logger, CustomResourceTemplate)
```

### Prompts

```go

var CustomPrompt = &Prompt{
	Name:        "custom_prompt",
	Description: "A custom prompt from " + util.AppName,
	Content:     `This is a prompt!`,
}

mcpServer, _ := mcpserver.NewServer(ctx, as, logger)
err := mcpServer.AddPrompts(as, logger, CustomPrompt)
```

### Web Admin Interface

Access the web admin interface at `/mcp` to:

- View registered tools
- Test tool execution
- Inspect request/response data
- Debug tool implementations

## Example Tools

The module includes example tools to demonstrate usage:

- **Random Number Generator** - Generates random numbers within specified ranges

## Configuration

The MCP module uses your application's standard configuration patterns:

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
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
