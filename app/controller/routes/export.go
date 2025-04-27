package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/cproject"
)

func exportRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/p/{key}/export", cproject.ProjectExportOverview)

	makeRoute(r, http.MethodGet, "/p/{key}/export/config", cproject.ProjectExportConfigForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/config", cproject.ProjectExportConfigSave)

	makeRoute(r, http.MethodGet, "/p/{key}/export/groups", cproject.ProjectExportGroupsEdit)
	makeRoute(r, http.MethodPost, "/p/{key}/export/groups", cproject.ProjectExportGroupsSave)

	makeRoute(r, http.MethodGet, "/p/{key}/export/jsonschema", cproject.ProjectExportJSONSchema)

	makeRoute(r, http.MethodGet, "/p/{key}/export/models/create/new", cproject.ProjectExportModelNew)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/create/new", cproject.ProjectExportModelCreate)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/create/derive", cproject.ProjectExportDeriveForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/create/derive", cproject.ProjectExportDerive)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}", cproject.ProjectExportModelDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/seeddata", cproject.ProjectExportModelSeedData)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/delete", cproject.ProjectExportModelDelete)

	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}", cproject.ProjectExportEnumDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/new", cproject.ProjectExportEnumNew)
	makeRoute(r, http.MethodPost, "/p/{key}/export/enums/new", cproject.ProjectExportEnumCreate)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}/edit", cproject.ProjectExportEnumForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/enums/{enum}/edit", cproject.ProjectExportEnumSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}/delete", cproject.ProjectExportEnumDelete)
}
