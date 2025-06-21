package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/clib"
)

func adminRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/admin", clib.Admin)
	makeRoute(r, http.MethodGet, "/admin/mcp", clib.MCPIndex)
	makeRoute(r, http.MethodGet, "/admin/mcp/sse", clib.MCPServe)
	makeRoute(r, http.MethodPost, "/admin/mcp/sse", clib.MCPServe)
	makeRoute(r, http.MethodGet, "/admin/mcp/sse/{path:.*}", clib.MCPServe)
	makeRoute(r, http.MethodPost, "/admin/mcp/sse/{path:.*}", clib.MCPServe)
	makeRoute(r, http.MethodGet, "/admin/mcp/tool/{tool}", clib.MCPTool)
	makeRoute(r, http.MethodPost, "/admin/mcp/tool/{tool}", clib.MCPToolRun)
	makeRoute(r, http.MethodGet, "/admin/task", clib.TaskList)
	makeRoute(r, http.MethodGet, "/admin/task/{key}", clib.TaskDetail)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/remove", clib.TaskRemove)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/run", clib.TaskRun)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/start", clib.TaskStart)
	makeRoute(r, http.MethodGet, "/admin/{path:.*}", clib.Admin)
	makeRoute(r, http.MethodPost, "/admin/{path:.*}", clib.Admin)
}
