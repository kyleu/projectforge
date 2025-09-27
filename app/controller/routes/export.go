package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/cproject/cexport"
)

func exportRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/p/{key}/export", cexport.ProjectExportOverview)
	makeRoute(r, http.MethodGet, "/p/{key}/export/jsonschema", cexport.ProjectExportJSONSchema)
	makeRoute(r, http.MethodGet, "/p/{key}/export/jsonschema/write", cexport.ProjectExportWriteJSONSchema)

	makeRoute(r, http.MethodGet, "/p/{key}/export/config", cexport.ProjectExportConfigForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/config", cexport.ProjectExportConfigSave)

	makeRoute(r, http.MethodGet, "/p/{key}/export/groups", cexport.ProjectExportGroupsEdit)
	makeRoute(r, http.MethodPost, "/p/{key}/export/groups", cexport.ProjectExportGroupsSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/new", cexport.ProjectExportEnumNew)
	makeRoute(r, http.MethodPost, "/p/{key}/export/enums/new", cexport.ProjectExportEnumCreate)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}", cexport.ProjectExportEnumDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}/jsonschema", cexport.ProjectExportEnumJSONSchema)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}/edit", cexport.ProjectExportEnumForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/enums/{enum}/edit", cexport.ProjectExportEnumSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/enums/{enum}/delete", cexport.ProjectExportEnumDelete)

	makeRoute(r, http.MethodGet, "/p/{key}/export/events/create/new", cexport.ProjectExportEventNew)
	makeRoute(r, http.MethodPost, "/p/{key}/export/events/create/new", cexport.ProjectExportEventCreate)
	makeRoute(r, http.MethodGet, "/p/{key}/export/events/create/derive", cexport.ProjectExportDeriveForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/events/create/derive", cexport.ProjectExportDerive)
	makeRoute(r, http.MethodGet, "/p/{key}/export/events/{event}", cexport.ProjectExportEventDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/export/events/{event}/edit", cexport.ProjectExportEventForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/events/{event}/edit", cexport.ProjectExportEventSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/events/{event}/delete", cexport.ProjectExportEventDelete)

	makeRoute(r, http.MethodGet, "/p/{key}/export/models/create/new", cexport.ProjectExportModelNew)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/create/new", cexport.ProjectExportModelCreate)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/create/derive", cexport.ProjectExportDeriveForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/create/derive", cexport.ProjectExportDerive)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}", cexport.ProjectExportModelDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/jsonschema", cexport.ProjectExportModelJSONSchema)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/seeddata", cexport.ProjectExportModelSeedData)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/edit", cexport.ProjectExportModelForm)
	makeRoute(r, http.MethodPost, "/p/{key}/export/models/{model}/edit", cexport.ProjectExportModelSave)
	makeRoute(r, http.MethodGet, "/p/{key}/export/models/{model}/delete", cexport.ProjectExportModelDelete)

	makeRoute(r, http.MethodGet, "/p/{key}/export/extra/{t}", cexport.ProjectExportExtraTypeDetail)
}
