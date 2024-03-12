package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func notebookRoutes(r *router.Router) {
	r.GET("/notebook", clib.Notebook)
	r.GET("/notebook/action/{act}", clib.NotebookAction)
}
