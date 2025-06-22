package mcpserver

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type Server struct {
	MCP       *server.MCPServer `json:"-"`
	State     *app.State        `json:"-"`
	Resources Resources         `json:"resources"`
	Prompts   Prompts           `json:"prompts"`
	Tools     Tools             `json:"tools"`
	HTTP      http.Handler      `json:"-"`
}

func NewServer(ctx context.Context, as *app.State, logger util.Logger) (*Server, error) {
	t := util.TimerStart()
	ms := server.NewMCPServer(util.AppName, as.BuildInfo.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)
	mcp := &Server{MCP: ms}
	// $PF_SECTION_START(tools)$
	// if err := mcp.AddTools(as, logger, ...); err != nil {
	// 	return nil, err
	// }
	// $PF_SECTION_END(tools)$
	logger.Debugf("MCP server initialized in [%s] with [%d] resources, [%d] tools, and [%d] prompts", t.EndString(), len(mcp.Resources), len(mcp.Tools), len(mcp.Prompts))
	return mcp, nil
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
