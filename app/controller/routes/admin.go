package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/clib"
)

func adminRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/admin", clib.Admin)
	makeRoute(r, http.MethodGet, "/admin/mcp", clib.MCPIndex)
	makeRoute(r, http.MethodPost, "/admin/mcp", clib.MCPTaskSSE)
	makeRoute(r, http.MethodGet, "/admin/mcp/tool/{tool}", clib.MCPTask)
	makeRoute(r, http.MethodPost, "/admin/mcp/tool/{tool}", clib.MCPTaskRun)
	makeRoute(r, http.MethodGet, "/admin/task", clib.TaskList)
	makeRoute(r, http.MethodGet, "/admin/task/{key}", clib.TaskDetail)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/remove", clib.TaskRemove)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/run", clib.TaskRun)
	makeRoute(r, http.MethodGet, "/admin/task/{key}/start", clib.TaskStart)
	makeRoute(r, http.MethodGet, "/admin/{path:.*}", clib.Admin)
	makeRoute(r, http.MethodPost, "/admin/{path:.*}", clib.Admin)
}
