package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/controller/cproject"
)

func projectRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/welcome", cproject.Welcome)
	makeRoute(r, http.MethodPost, "/welcome", cproject.WelcomeLoad)
	makeRoute(r, http.MethodGet, "/welcome/changedir", cproject.ChangeDir)
	makeRoute(r, http.MethodGet, "/welcome/process", cproject.WelcomeProcess)

	makeRoute(r, http.MethodGet, "/doctor", cproject.Doctor)
	makeRoute(r, http.MethodGet, "/doctor/all", cproject.DoctorRunAll)
	makeRoute(r, http.MethodGet, "/doctor/{key}", cproject.DoctorRun)
	makeRoute(r, http.MethodGet, "/doctor/{key}/solve", cproject.DoctorSolve)

	makeRoute(r, http.MethodGet, "/p", cproject.ProjectList)
	makeRoute(r, http.MethodGet, "/p/search", cproject.ProjectSearchAll)
	makeRoute(r, http.MethodGet, "/p/new", cproject.ProjectForm)
	makeRoute(r, http.MethodPost, "/p/new", cproject.ProjectCreate)
	makeRoute(r, http.MethodGet, "/p/{key}", cproject.ProjectDetail)
	makeRoute(r, http.MethodGet, "/p/{key}/edit", cproject.ProjectEdit)
	makeRoute(r, http.MethodPost, "/p/{key}/edit", cproject.ProjectSave)
	makeRoute(r, http.MethodGet, "/p/{key}/stats", cproject.ProjectFileStats)
	makeRoute(r, http.MethodGet, "/p/{key}/start", cproject.ProjectStart)

	exportRoutes(r)

	makeRoute(r, http.MethodGet, "/p/{key}/fs", cproject.ProjectFileRoot)
	makeRoute(r, http.MethodGet, "/p/{key}/fs/{path:.*}", cproject.ProjectFile)
	makeRoute(r, http.MethodGet, "/p/{key}/search", cproject.ProjectSearch)

	makeRoute(r, http.MethodGet, "/p/{key}/palette/{palette}", cproject.ProjectThemePalette)
	makeRoute(r, http.MethodGet, "/p/{key}/palette/{palette}/{theme}", cproject.ProjectThemeSave)

	makeRoute(r, http.MethodGet, "/svg/{key}", cproject.SVGList)
	makeRoute(r, http.MethodGet, "/svg/{key}/x/add", cproject.SVGAdd)
	makeRoute(r, http.MethodGet, "/svg/{key}/x/build", cproject.SVGBuild)
	makeRoute(r, http.MethodGet, "/svg/{key}/x/refreshapp", cproject.SVGRefreshApp)
	makeRoute(r, http.MethodGet, "/svg/{key}/{icon}", cproject.SVGDetail)
	makeRoute(r, http.MethodGet, "/svg/{key}/{icon}/setapp", cproject.SVGSetApp)
	makeRoute(r, http.MethodGet, "/svg/{key}/{icon}/remove", cproject.SVGRemove)

	makeRoute(r, http.MethodGet, "/git", cproject.GitActionAll)
	makeRoute(r, http.MethodGet, "/git/all/{act}", cproject.GitActionAll)

	makeRoute(r, http.MethodGet, "/git/{key}", cproject.GitAction)
	makeRoute(r, http.MethodGet, "/git/{key}/{act}", cproject.GitAction)

	makeRoute(r, http.MethodGet, "/run/{act}", cproject.RunAllActions)
	makeRoute(r, http.MethodGet, "/run/{key}/{act}", cproject.RunAction)

	makeRoute(r, http.MethodGet, "/test", cproject.TestList)
	makeRoute(r, http.MethodGet, "/test/{key}", cproject.TestRun)
}
