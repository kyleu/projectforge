package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var (
	buildInfo *app.BuildInfo
	debug     bool
)

func InitMCP(bi *app.BuildInfo, dbg bool) {
	buildInfo = bi
	debug = dbg
}

type Server struct {
	MCP   *server.MCPServer
	Tools map[string]*Tool
}

func NewServer(ctx context.Context, version string) (*Server, error) {
	ms := server.NewMCPServer(util.AppName, version)
	mcp := &Server{MCP: ms, Tools: make(map[string]*Tool)}
	// $PF_SECTION_START(tools)$
	// $PF_SECTION_END(tools)$
	return mcp, nil
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

func MCPConfig() (*app.BuildInfo, bool) {
	return buildInfo, debug
}
