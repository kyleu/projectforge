package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func notebookRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/notebook", clib.Notebook)
	makeRoute(r, http.MethodGet, "/notebook/view/{path:.*}", clib.Notebook)
	makeRoute(r, http.MethodGet, "/notebook/files", clib.NotebookFiles)
	makeRoute(r, http.MethodGet, "/notebook/files/{path:.*}", clib.NotebookFiles)
	makeRoute(r, http.MethodGet, "/notebook/edit/{path:.*}", clib.NotebookFileEdit)
	makeRoute(r, http.MethodPost, "/notebook/edit/{path:.*}", clib.NotebookFileSave)
	makeRoute(r, http.MethodGet, "/notebook/action/{act}", clib.NotebookAction)
}
