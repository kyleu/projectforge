package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func scriptingRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/admin/scripting", clib.ScriptingList)
	makeRoute(r, http.MethodGet, "/admin/scripting/new", clib.ScriptingNew)
	makeRoute(r, http.MethodPost, "/admin/scripting/new", clib.ScriptingCreate)
	makeRoute(r, http.MethodGet, "/admin/scripting/{key}", clib.ScriptingDetail)
	makeRoute(r, http.MethodGet, "/admin/scripting/{key}/edit", clib.ScriptingForm)
	makeRoute(r, http.MethodPost, "/admin/scripting/{key}/edit", clib.ScriptingSave)
	makeRoute(r, http.MethodGet, "/admin/scripting/{key}/delete", clib.ScriptingDelete)
}
