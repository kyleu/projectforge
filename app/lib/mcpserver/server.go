package mcpserver

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type Server struct {
	MCP   *server.MCPServer `json:"-"`
	State *app.State        `json:"-"`
	Tools Tools
	HTTP  http.Handler `json:"-"`
}

func NewServer(ctx context.Context, as *app.State, logger util.Logger, tools ...*Tool) (*Server, error) {
	ms := server.NewMCPServer(util.AppName, as.BuildInfo.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)
	mcp := &Server{MCP: ms}
	// $PF_SECTION_START(tools)$
	if err := mcp.AddTools(as, logger, AllTools...); err != nil {
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
	logger.Debugf("MCP HTTP request: %s %s", r.Method, r.URL.Path)
	if s.HTTP == nil {
		s.HTTP = server.NewStreamableHTTPServer(s.MCP)
	}
	s.HTTP.ServeHTTP(w, r)
}
