package routes

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/clib"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

func AppRoutes(as *app.State, logger util.Logger) fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", controller.Home)
	r.GET("/healthcheck", clib.Healthcheck)
	r.GET("/about", clib.About)
	r.GET("/theme", clib.ThemeList)
	r.GET("/theme/{key}", clib.ThemeEdit)
	r.POST("/theme/{key}", clib.ThemeSave){{{ if .HasModule "search" }}}
	r.GET(cutil.DefaultSearchPath, clib.Search){{{ end }}}

	r.GET(cutil.DefaultProfilePath, clib.Profile)
	r.POST(cutil.DefaultProfilePath, clib.ProfileSave){{{ if .HasModule "oauth" }}}
	r.GET("/auth/{key}", clib.AuthDetail)
	r.GET("/auth/callback/{key}", clib.AuthCallback)
	r.GET("/auth/logout/{key}", clib.AuthLogout){{{ end }}}{{{ if.HasModule "export" }}}

	generatedRoutes(r){{{ end }}}

	// $PF_SECTION_START(routes)$
	// $PF_SECTION_END(routes)${{{ if .HasModule "docbrowse" }}}

	r.GET("/docs", clib.Docs)
	r.GET("/docs/{path:*}", clib.Docs){{{ end }}}{{{ if .HasModule "graphql" }}}

	r.GET("/graphql", controller.GraphQLIndex)
	r.GET("/graphql/{key}", controller.GraphQLDetail)
	r.POST("/graphql/{key}", controller.GraphQLRun){{{ end }}}

	r.GET("/admin", clib.Admin){{{ if.HasModule "audit" }}}
	r.GET("/admin/audit", clib.AuditList)
	r.GET("/admin/audit/random", clib.AuditCreateFormRandom)
	r.GET("/admin/audit/new", clib.AuditCreateForm)
	r.POST("/admin/audit/new", clib.AuditCreate)
	r.GET("/admin/audit/record/{id}", clib.RecordDetail)
	r.GET("/admin/audit/{id}", clib.AuditDetail)
	r.GET("/admin/audit/{id}/edit", clib.AuditEditForm)
	r.POST("/admin/audit/{id}/edit", clib.AuditEdit)
	r.GET("/admin/audit/{id}/delete", clib.AuditDelete){{{ end }}}{{{ if .HasModule "databaseui" }}}
	r.GET("/admin/database", clib.DatabaseList)
	r.GET("/admin/database/{key}", clib.DatabaseDetail)
	r.GET("/admin/database/{key}/{act}", clib.DatabaseAction)
	r.GET("/admin/database/{key}/tables/{schema}/{table}", clib.DatabaseTableView){{{ if .DatabaseUISQLEditor }}}
	r.POST("/admin/database/{key}/sql", clib.DatabaseSQLRun){{{ end }}}{{{ end }}}{{{ if .HasModule "sandbox" }}}
	r.GET("/admin/sandbox", controller.SandboxList)
	r.GET("/admin/sandbox/{key}", controller.SandboxRun){{{ end }}}
	r.GET("/admin/{path:*}", clib.Admin)
	r.POST("/admin/{path:*}", clib.Admin)

	r.GET("/favicon.ico", clib.Favicon)
	r.GET("/robots.txt", clib.RobotsTxt)
	r.GET("/assets/{_:*}", clib.Static)

	r.OPTIONS("/", controller.Options)
	r.OPTIONS("/{_:*}", controller.Options)
	r.NotFound = controller.NotFound

	clib.AppRoutesList = r.List()

	p := httpmetrics.NewMetrics(util.AppKey, logger)
	return fasthttp.CompressHandlerLevel(p.WrapHandler(r), fasthttp.CompressBestSpeed)
}
