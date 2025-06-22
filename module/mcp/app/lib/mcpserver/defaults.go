package mcpserver

import (
	"context"
	"fmt"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var defaultServer *Server
var defaultTools Tools

func AddDefaultTools(ts ...*Tool) {
	defaultTools = append(defaultTools, ts...)
}

func ClearDefaultTools() {
	defaultTools = nil
}

func GetDefaultServer(ctx context.Context, as *app.State, logger util.Logger) (*Server, error) {
	if defaultServer != nil {
		return defaultServer, nil
	}
	ret, err := NewServer(ctx, as, logger)
	if err != nil {
		return nil, err
	}
	defaultServer = ret
	return defaultServer, nil
}

func CurrentDefaultServer() *Server {
	return defaultServer
}

const usageCLI = `{
  "mcpServers": {
    "%s": {
      "command": "%s",
      "args": ["mcp"]
    }
  }
}`

func UsageCLI() string {
	return fmt.Sprintf(usageCLI, util.AppCmd, util.AppCmd)
}

const usageHTTP = `{
  "mcpServers": {
    "%s": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "http://localhost:%d/admin/mcp"]
    }
  }
}`

func UsageHTTP() string {
	return fmt.Sprintf(usageHTTP, util.AppCmd, util.AppPort)
}
