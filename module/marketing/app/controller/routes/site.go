package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/clib"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func SiteRoutes(logger util.Logger) (http.Handler, error) {
	r := mux.NewRouter()

	makeRoute(r, http.MethodGet, cutil.DefaultProfilePath, clib.ProfileSite)
	makeRoute(r, http.MethodPost, cutil.DefaultProfilePath, clib.ProfileSave){{{ if .HasAccount }}}
	makeRoute(r, http.MethodGet, "/auth/{key}", clib.AuthDetail)
	makeRoute(r, http.MethodGet, "/auth/callback/{key}", clib.AuthCallback)
	makeRoute(r, http.MethodGet, "/auth/logout/{key}", clib.AuthLogout){{{ end }}}

	makeRoute(r, http.MethodGet, "/favicon.ico", clib.Favicon)
	makeRoute(r, http.MethodGet, "/assets/{path:.*}", clib.Static)

	makeRoute(r, http.MethodGet, "/", controller.Site)
	makeRoute(r, http.MethodGet, "/{path:.*}", controller.Site)

	r.Methods(http.MethodOptions).HandlerFunc(controller.Options)
	r.Methods(http.MethodHead).HandlerFunc(controller.Head)

	return cutil.WireRouter(r, controller.NotFoundAction, logger)
}
