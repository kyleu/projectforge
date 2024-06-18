package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/clib"
)

func themeRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/theme", clib.ThemeList)
	makeRoute(r, http.MethodGet, "/theme/{key}", clib.ThemeEdit)
	makeRoute(r, http.MethodPost, "/theme/{key}", clib.ThemeSave)
	makeRoute(r, http.MethodGet, "/theme/{key}/remove", clib.ThemeRemove)
}
