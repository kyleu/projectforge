package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/clib"
)

func mcpRoutes(base string, r *mux.Router) {
	if base == "" {
		base = "/mcp"
	}
	makeRoute(r, http.MethodGet, base, clib.MCPIndex)
	makeRoute(r, http.MethodPost, base, clib.MCPServe)
	makeRoute(r, http.MethodGet, base+"/tool/{tool}", clib.MCPTool)
	makeRoute(r, http.MethodPost, base+"/tool/{tool}", clib.MCPToolRun)
}
