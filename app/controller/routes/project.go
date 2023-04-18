package routes

import (
	"github.com/fasthttp/router"

	"projectforge.dev/projectforge/app/controller/cproject"
)

func projectRoutes(r *router.Router) {
	r.GET("/welcome", cproject.Welcome)
	r.POST("/welcome", cproject.WelcomeResult)

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
	r.GET("/p/{key}/start", cproject.ProjectStart)

	exportRoutes(r)

	r.GET("/p/{key}/fs", cproject.ProjectFileRoot)
	r.GET("/p/{key}/fs/{path:*}", cproject.ProjectFile)
	r.GET("/p/{key}/search", cproject.ProjectSearch)

	r.GET("/svg/{key}", cproject.SVGList)
	r.GET("/svg/{key}/x/add", cproject.SVGAdd)
	r.GET("/svg/{key}/x/build", cproject.SVGBuild)
	r.GET("/svg/{key}/x/refreshapp", cproject.SVGRefreshApp)
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
