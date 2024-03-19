package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func execRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/admin/exec", clib.ExecList)
	makeRoute(r, http.MethodGet, "/admin/exec/new", clib.ExecForm)
	makeRoute(r, http.MethodPost, "/admin/exec/new", clib.ExecNew)
	makeRoute(r, http.MethodGet, "/admin/exec/{key}/{idx}", clib.ExecDetail){{{ if .HasModule "websocket" }}}
	makeRoute(r, http.MethodGet, "/admin/exec/{key}/{idx}/connect", clib.ExecSocket){{{ end }}}
	makeRoute(r, http.MethodGet, "/admin/exec/{key}/{idx}/kill", clib.ExecKill)
}
