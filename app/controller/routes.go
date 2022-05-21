// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/lib/telemetry/httpmetrics"
	"projectforge.dev/projectforge/app/util"
)

//nolint
func AppRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", Home)
	r.GET("/healthcheck", Healthcheck)
	r.GET("/about", About)
	r.GET("/theme", ThemeList)
	r.GET("/theme/{key}", ThemeEdit)
	r.POST("/theme/{key}", ThemeSave)
	r.GET(defaultSearchPath, Search)

	r.GET(defaultProfilePath, Profile)
	r.POST(defaultProfilePath, ProfileSave)

	// $PF_SECTION_START(routes)$
	r.GET("/welcome", Welcome)
	r.POST("/welcome", WelcomeResult)

	r.GET("/doctor", Doctor)
	r.GET("/doctor/all", DoctorRunAll)
	r.GET("/doctor/{key}", DoctorRun)

	r.GET("/p", ProjectList)
	r.GET("/p/search", ProjectSearchAll)
	r.GET("/p/new", ProjectForm)
	r.POST("/p/new", ProjectCreate)
	r.GET("/p/{key}", ProjectDetail)
	r.GET("/p/{key}/edit", ProjectEdit)
	r.POST("/p/{key}/edit", ProjectSave)
	r.GET("/p/{key}/fs", ProjectFileRoot)
	r.GET("/p/{key}/stats", ProjectFileStats)
	r.GET("/p/{key}/export/{model}", ProjectExportModelDetail)
	r.GET("/p/{key}/fs/{path:*}", ProjectFile)
	r.GET("/p/{key}/search", ProjectSearch)

	r.GET("/svg/{key}", SVGList)
	r.GET("/svg/{key}/x/add", SVGAdd)
	r.GET("/svg/{key}/x/build", SVGBuild)
	r.GET("/svg/{key}/{icon}", SVGDetail)
	r.GET("/svg/{key}/{icon}/setapp", SVGSetApp)
	r.GET("/svg/{key}/{icon}/remove", SVGRemove)

	r.GET("/git", GitActionAll)
	r.GET("/git/all/{act}", GitActionAll)

	r.GET("/git/{key}", GitAction)
	r.GET("/git/{key}/{act}", GitAction)

	r.GET("/run/{act}", RunAllActions)
	r.GET("/run/{key}/{act}", RunAction)

	r.GET("/m", ModuleList)
	r.GET("/m/{key}", ModuleDetail)
	r.GET("/m/{key}/fs", ModuleFileRoot)
	r.GET("/m/{key}/fs/{path:*}", ModuleFile)
	r.GET("/m/{key}/search", ModuleSearch)

	r.GET("/test", TestList)
	r.GET("/test/{key}", TestRun)

	r.GET("/theme/palette/{palette}", ThemePalette)
	r.GET("/theme/color/{color}", ThemeColor)
	r.GET("/theme/preview/{palette}/{key}", ThemePreview)
	// $PF_SECTION_END(routes)$

	r.GET("/docs", Docs)
	r.GET("/docs/{path:*}", Docs)

	r.GET("/admin", Admin)
	r.GET("/admin/{path:*}", Admin)

	r.GET("/favicon.ico", Favicon)
	r.GET("/robots.txt", RobotsTxt)
	r.GET("/assets/{_:*}", Static)

	r.OPTIONS("/", Options)
	r.OPTIONS("/{_:*}", Options)
	r.NotFound = NotFound

	p := httpmetrics.NewMetrics(util.AppKey)
	return fasthttp.CompressHandlerLevel(p.WrapHandler(r), fasthttp.CompressBestSpeed)
}
