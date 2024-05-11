package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func proxyRoutes(pth string, r *mux.Router) {
	makeRoute(r, http.MethodGet, pth, clib.ProxyIndex)
	makeRoute(r, http.MethodGet, pth+"/{svc}", clib.ProxyHandle)
	makeRoute(r, http.MethodGet, pth+"/{svc}/{path:.*}", clib.ProxyHandle)
}
