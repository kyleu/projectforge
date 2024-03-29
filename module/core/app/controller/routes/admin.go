package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/clib"
)

func adminRoutes(r *mux.Router) {
	makeRoute(r, http.MethodGet, "/admin", clib.Admin){{{ if.HasModule "audit" }}}
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
	makeRoute(r, http.MethodPost, "/admin/database/{key}/sql", clib.DatabaseSQLRun){{{ end }}}{{{ end }}}{{{ if .HasModule "schedule" }}}
	makeRoute(r, http.MethodGet, "/admin/schedule", clib.ScheduleList)
	makeRoute(r, http.MethodGet, "/admin/schedule/{id}", clib.ScheduleDetail){{{ end }}}{{{ if .HasModule "queue" }}}
	makeRoute(r, http.MethodGet, "/admin/queue", clib.QueueIndex)
	makeRoute(r, http.MethodPost, "/admin/queue", clib.QueueSend){{{ end }}}{{{ if .HasModule "sandbox" }}}
	makeRoute(r, http.MethodGet, "/admin/sandbox", clib.SandboxList)
	makeRoute(r, http.MethodGet, "/admin/sandbox/{key}", clib.SandboxRun){{{ end }}}
	makeRoute(r, http.MethodGet, "/admin/{path:.*}", clib.Admin)
	makeRoute(r, http.MethodPost, "/admin/{path:.*}", clib.Admin)
}
