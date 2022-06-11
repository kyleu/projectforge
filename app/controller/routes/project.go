package routes

import (
	"github.com/fasthttp/router"

	"projectforge.dev/projectforge/app/controller/cproject"
)

func projectRoutes(r *router.Router) {
	r.GET("/doctor", cproject.Doctor)
	r.GET("/doctor/all", cproject.DoctorRunAll)
	r.GET("/doctor/{key}", cproject.DoctorRun)

	r.GET("/p", cproject.ProjectList)
	r.GET("/p/search", cproject.ProjectSearchAll)
	r.GET("/p/new", cproject.ProjectForm)
	r.POST("/p/new", cproject.ProjectCreate)
	r.GET("/p/{key}", cproject.ProjectDetail)
	r.GET("/p/{key}/edit", cproject.ProjectEdit)
	r.POST("/p/{key}/edit", cproject.ProjectSave)
	r.GET("/p/{key}/stats", cproject.ProjectFileStats)
	r.GET("/p/{key}/export", cproject.ProjectExportOverview)
	r.GET("/p/{key}/export/config", cproject.ProjectExportConfigForm)
	r.POST("/p/{key}/export/config", cproject.ProjectExportConfigSave)
	r.GET("/p/{key}/export/models/{model}", cproject.ProjectExportModelDetail)
	r.GET("/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelForm)
	r.POST("/p/{key}/export/models/{model}/edit", cproject.ProjectExportModelSave)
	r.GET("/p/{key}/export/models/{model}/delete", cproject.ProjectExportModelDelete)
	r.GET("/p/{key}/fs", cproject.ProjectFileRoot)
	r.GET("/p/{key}/fs/{path:*}", cproject.ProjectFile)
	r.GET("/p/{key}/search", cproject.ProjectSearch)

	r.GET("/svg/{key}", cproject.SVGList)
	r.GET("/svg/{key}/x/add", cproject.SVGAdd)
	r.GET("/svg/{key}/x/build", cproject.SVGBuild)
	r.GET("/svg/{key}/{icon}", cproject.SVGDetail)
	r.GET("/svg/{key}/{icon}/setapp", cproject.SVGSetApp)
	r.GET("/svg/{key}/{icon}/remove", cproject.SVGRemove)

	r.GET("/git", cproject.GitActionAll)
	r.GET("/git/all/{act}", cproject.GitActionAll)

	r.GET("/git/{key}", cproject.GitAction)
	r.GET("/git/{key}/{act}", cproject.GitAction)

	r.GET("/run/{act}", cproject.RunAllActions)
	r.GET("/run/{key}/{act}", cproject.RunAction)

	r.GET("/test", cproject.TestList)
	r.GET("/test/{key}", cproject.TestRun)
}
