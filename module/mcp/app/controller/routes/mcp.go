package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib/cmcp"
)

func mcpRoutes(base string, r *mux.Router) {
	if base == "" {
		base = "/mcp"
	}
	makeRoute(r, http.MethodGet, base, cmcp.MCPIndex)
	makeRoute(r, http.MethodPost, base, cmcp.MCPServe)
	makeRoute(r, http.MethodGet, base+"/resource/{resource}", cmcp.MCPResource)
	makeRoute(r, http.MethodGet, base+"/tool/{tool}", cmcp.MCPTool)
	makeRoute(r, http.MethodPost, base+"/tool/{tool}", cmcp.MCPToolRun)
	makeRoute(r, http.MethodGet, base+"/prompt/{prompt}", cmcp.MCPPrompt)
	makeRoute(r, http.MethodPost, base+"/prompt/{prompt}", cmcp.MCPPromptRun)
}
