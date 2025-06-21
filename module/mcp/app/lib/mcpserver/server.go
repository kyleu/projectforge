package mcpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type Server struct {
	MCP   *server.MCPServer `json:"-"`
	State *app.State        `json:"-"`
	Tools Tools
	HTTP  http.Handler `json:"-"`
}

func NewServer(ctx context.Context, as *app.State, logger util.Logger, tools ...*Tool) (*Server, error) {
	ms := server.NewMCPServer(util.AppName, as.BuildInfo.Version)
	mcp := &Server{MCP: ms}
	// $PF_SECTION_START(tools)$
	if err := mcp.AddTools(as, logger, ProjectsTool); err != nil {
		return nil, err
	}
	// $PF_SECTION_END(tools)$
	if err := mcp.AddTools(as, logger, tools...); err != nil {
		return nil, err
	}
	return mcp, nil
}

func (s *Server) AddTools(as *app.State, logger util.Logger, tools ...*Tool) error {
	for _, tl := range tools {
		s.Tools = append(s.Tools, tl)
		m, err := tl.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddTool(m, tl.Handler(as, logger))
	}
	return nil
}

func (s *Server) ServeCLI(ctx context.Context) error {
	return server.ServeStdio(s.MCP)
}

func (s *Server) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request, logger util.Logger) {
	logger.Debugf("MCP SSE request: %s %s", r.Method, r.URL.Path)
	if s.HTTP == nil {
		s.HTTP = server.NewSSEServer(s.MCP, server.WithBaseURL("/admin/mcp/sse")).SSEHandler()
	}
	s.HTTP.ServeHTTP(w, r)
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
      "args": ["-y", "mcp-remote", "http://localhost:%d/admin/mcp/sse"]
    }
  }
}`

func UsageHTTP() string {
	return fmt.Sprintf(usageHTTP, util.AppCmd, util.AppPort)
}
