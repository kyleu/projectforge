package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/cmodule"
)

func moduleRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/m", cmodule.ModuleList)
	makeRoute(r, http.MethodGet, "/m/{key}", cmodule.ModuleDetail)
	makeRoute(r, http.MethodGet, "/m/{key}/fs", cmodule.ModuleFileRoot)
	makeRoute(r, http.MethodGet, "/m/{key}/fs/{path:.*}", cmodule.ModuleFile)
	makeRoute(r, http.MethodGet, "/m/{key}/search", cmodule.ModuleSearch)
}
