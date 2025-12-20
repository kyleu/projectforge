package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/ctest"
)

func testRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/test", ctest.TestList)
	makeRoute(r, http.MethodGet, "/test/{key}", ctest.TestRun)
	makeRoute(r, http.MethodGet, "/test/jsonschema/{filename}", ctest.JSONSchemaTestFile)

	makeRoute(r, http.MethodGet, "/testbed", controller.Testbed)
	makeRoute(r, http.MethodPost, "/testbed", controller.Testbed)
}
