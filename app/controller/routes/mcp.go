package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/clib/cmcp"
)

func mcpRoutes(base string, r *mux.Router) {
	if base == "" {
		base = "/mcp"
	}
	makeRoute(r, http.MethodGet, base, cmcp.MCPIndex)
	makeRoute(r, http.MethodPost, base, cmcp.MCPServe)
	makeRoute(r, http.MethodGet, base+"/resource/{resource}", cmcp.MCPResource)
	makeRoute(r, http.MethodGet, base+"/resourcetemplate/{rt}", cmcp.MCPResourceTemplate)
	makeRoute(r, http.MethodPost, base+"/resourcetemplate/{rt}", cmcp.MCPResourceTemplateRun)
	makeRoute(r, http.MethodGet, base+"/tool/{tool}", cmcp.MCPTool)
	makeRoute(r, http.MethodPost, base+"/tool/{tool}", cmcp.MCPToolRun)
	makeRoute(r, http.MethodGet, base+"/prompt/{prompt}", cmcp.MCPPrompt)
}
