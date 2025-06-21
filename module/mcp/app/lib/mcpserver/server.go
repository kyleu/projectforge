package mcpserver

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type Server struct {
	MCP   *server.MCPServer `json:"-"`
	State *app.State        `json:"-"`
	Tools map[string]*Tool
	HTTP  http.Handler `json:"-"`
}

func NewServer(ctx context.Context, as *app.State, logger util.Logger, tools ...*Tool) (*Server, error) {
	ms := server.NewMCPServer(util.AppName, as.BuildInfo.Version)
	mcp := &Server{MCP: ms, Tools: make(map[string]*Tool)}
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
	for _, tool := range tools {
		s.Tools[tool.Name] = tool
		t, err := tool.ToMCP()
		if err != nil {
			return err
		}
		s.MCP.AddTool(t, tool.Handler(as, logger))
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
