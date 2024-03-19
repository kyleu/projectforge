package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func harRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/har", clib.HarList)
	makeRoute(r, http.MethodPost, "/har", clib.HarUpload)
	makeRoute(r, http.MethodGet, "/har/{key}", clib.HarDetail)
	makeRoute(r, http.MethodGet, "/har/{key}/delete", clib.HarDelete)
	makeRoute(r, http.MethodGet, "/har/{key}/trim", clib.HarTrim)
}
