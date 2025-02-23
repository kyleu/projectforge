package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var (
	BuildInfo *app.BuildInfo
	ConfigDir string
	Debug     bool
)

type Server struct {
	MCP   *server.MCPServer
	Tools map[string]*Tool
}

func NewServer(ctx context.Context, version string) (*Server, error) {
	ms := server.NewMCPServer(util.AppName, version)
	return &Server{MCP: ms, Tools: make(map[string]*Tool)}, nil
}

func (s *Server) AddTools(tools ...*Tool) error {
	for _, tool := range tools {
		s.Tools[tool.Name] = tool
		t, err := tool.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddTool(t, tool.Handler())
	}
	return nil
}

func (s *Server) Serve(ctx context.Context) error {
	return server.ServeStdio(s.MCP)
}
