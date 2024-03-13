package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func notebookRoutes(r *router.Router) {
	r.GET("/notebook", clib.Notebook)
	r.GET("/notebook/view/{path:*}", clib.Notebook)
	r.GET("/notebook/files", clib.NotebookFiles)
	r.GET("/notebook/files/{path:*}", clib.NotebookFiles)
	r.GET("/notebook/edit/{path:*}", clib.NotebookFileEdit)
	r.POST("/notebook/edit/{path:*}", clib.NotebookFileSave)
	r.GET("/notebook/action/{act}", clib.NotebookAction)
}
