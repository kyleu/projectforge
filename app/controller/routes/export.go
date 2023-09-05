package routes

import (
	"github.com/fasthttp/router"

	"projectforge.dev/projectforge/app/controller/cproject"
)

func exportRoutes(r *router.Router) {
	r.GET("/p/{key}/export", cproject.ProjectExportOverview)

	r.GET("/p/{key}/export/config", cproject.ProjectExportConfigForm)
	r.POST("/p/{key}/export/config", cproject.ProjectExportConfigSave)

	r.GET("/p/{key}/export/groups", cproject.ProjectExportGroupsEdit)
	r.POST("/p/{key}/export/groups", cproject.ProjectExportGroupsSave)

	r.GET("/p/{key}/export/models/create/new", cproject.ProjectExportModelNew)
	r.POST("/p/{key}/export/models/create/new", cproject.ProjectExportModelCreate)
	r.GET("/p/{key}/export/models/create/derive", cproject.ProjectExportDeriveForm)
	r.POST("/p/{key}/export/models/create/derive", cproject.ProjectExportDerive)
	r.GET("/p/{key}/export/models/{model}", cproject.ProjectExportModelDetail)
	r.GET("/p/{key}/export/models/{model}/seeddata", cproject.ProjectExportModelSeedData)
	r.GET("/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelForm)
	r.POST("/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelSave)
	r.GET("/p/{key}/export/models/{model}/delete", cproject.ProjectExportModelDelete)

	r.GET("/p/{key}/export/enums/{enum}", cproject.ProjectExportEnumDetail)
	r.GET("/p/{key}/export/enums/{enum}/edit", cproject.ProjectExportEnumForm)
	r.POST("/p/{key}/export/enums/{enum}/edit", cproject.ProjectExportEnumSave)
	r.GET("/p/{key}/export/enums/{enum}/delete", cproject.ProjectExportEnumDelete)
}
