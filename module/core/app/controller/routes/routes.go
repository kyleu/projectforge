package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/clib"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func makeRoute(x *mux.Router, method string, path string, f http.HandlerFunc) {
	cutil.AddRoute(method, path)
	x.HandleFunc(path, f).Methods(method)
}

//nolint:revive
func AppRoutes(as *app.State, logger util.Logger) (http.Handler, error) {
	r := mux.NewRouter()

	makeRoute(r, http.MethodGet, "/", controller.Home)
	makeRoute(r, http.MethodGet, "/healthcheck", clib.Healthcheck)
	makeRoute(r, http.MethodGet, "/about", clib.About)

	makeRoute(r, http.MethodGet, cutil.DefaultProfilePath, clib.Profile)
	makeRoute(r, http.MethodPost, cutil.DefaultProfilePath, clib.ProfileSave){{{ if .HasAccount }}}
	makeRoute(r, http.MethodGet, "/auth/{key}", clib.AuthDetail)
	makeRoute(r, http.MethodGet, "/auth/callback/{key}", clib.AuthCallback)
	makeRoute(r, http.MethodGet, "/auth/logout/{key}", clib.AuthLogout){{{ end }}}{{{ if .HasModule "search" }}}
	makeRoute(r, http.MethodGet, cutil.DefaultSearchPath, clib.Search){{{ end }}}

	themeRoutes(r){{{ if.HasModule "export" }}}
	generatedRoutes(r){{{ end }}}

	// $PF_SECTION_START(routes)$
	// Add your custom routes here
	// $PF_SECTION_END(routes)${{{ if .HasModule "docbrowse" }}}

	makeRoute(r, http.MethodGet, "/docs", clib.Docs)
	makeRoute(r, http.MethodGet, "/docs/{path:.*}", clib.Docs){{{ end }}}{{{ if .HasModule "graphql" }}}

	makeRoute(r, http.MethodGet, "/graphql", controller.GraphQLIndex)
	makeRoute(r, http.MethodGet, "/graphql/{key}", controller.GraphQLDetail)
	makeRoute(r, http.MethodPost, "/graphql/{key}", controller.GraphQLRun){{{ end }}}

	makeRoute(r, http.MethodGet, "/admin", clib.Admin)
	makeRoute(r, http.MethodGet, "/admin/", clib.Admin)
	makeRoute(r, http.MethodPost, "/admin/", clib.Admin){{{ if.HasModule "audit" }}}
	makeRoute(r, http.MethodGet, "/admin/audit", clib.AuditList)
	makeRoute(r, http.MethodGet, "/admin/audit/random", clib.AuditCreateFormRandom)
	makeRoute(r, http.MethodGet, "/admin/audit/new", clib.AuditCreateForm)
	makeRoute(r, http.MethodPost, "/admin/audit/new", clib.AuditCreate)
	makeRoute(r, http.MethodGet, "/admin/audit/record/{id}/view", clib.RecordDetail)
	makeRoute(r, http.MethodGet, "/admin/audit/{id}", clib.AuditDetail)
	makeRoute(r, http.MethodGet, "/admin/audit/{id}/edit", clib.AuditEditForm)
	makeRoute(r, http.MethodPost, "/admin/audit/{id}/edit", clib.AuditEdit)
	makeRoute(r, http.MethodGet, "/admin/audit/{id}/delete", clib.AuditDelete){{{ end }}}{{{ if .HasModule "databaseui" }}}
	makeRoute(r, http.MethodGet, "/admin/database", clib.DatabaseList)
	makeRoute(r, http.MethodGet, "/admin/database/{key}", clib.DatabaseDetail)
	makeRoute(r, http.MethodGet, "/admin/database/{key}/{act}", clib.DatabaseAction)
	makeRoute(r, http.MethodGet, "/admin/database/{key}/tables/{schema}/{table}", clib.DatabaseTableView){{{ if .DatabaseUISQLEditor }}}
	makeRoute(r, http.MethodPost, "/admin/database/{key}/sql", clib.DatabaseSQLRun){{{ end }}}{{{ end }}}{{{ if .HasModule "sandbox" }}}
	makeRoute(r, http.MethodGet, "/admin/sandbox", controller.SandboxList)
	makeRoute(r, http.MethodGet, "/admin/sandbox/{key}", controller.SandboxRun){{{ end }}}{{{ if .HasModule "process" }}}
	execRoutes(r){{{ end }}}{{{ if .HasModule "scripting" }}}
	scriptingRoutes(r){{{ end }}}

	makeRoute(r, http.MethodGet, "/favicon.ico", clib.Favicon)
	makeRoute(r, http.MethodGet, "/robots.txt", clib.RobotsTxt)
	makeRoute(r, http.MethodGet, "/assets/{path:.*}", clib.Static)

	makeRoute(r, http.MethodOptions, "/", controller.Options)
	r.HandleFunc("/", controller.NotFoundAction)

	return cutil.WireRouter(r, logger)
}
