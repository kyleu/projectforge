package clib

import (
	"fmt"
	"net/http"
	"runtime"{{{ if .DangerousOK }}}
	"runtime/pprof"{{{ end }}}
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database/migrate"{{{ end }}}
	"{{{ .Package }}}/app/lib/log"{{{ if .HasAccount }}}
	"{{{ .Package }}}/app/lib/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vadmin"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	path := util.StringSplitAndTrim(strings.TrimPrefix(r.URL.Path, "/admin"), "/")
	key := "admin"
	if len(path) > 0 {
		key += "." + strings.Join(path, ".")
	}
	controller.Act(key, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		if len(path) == 0 {
			ps.SetTitleAndData("Administration", "administration")
			return controller.Render(w, r, as, {{{ if .HasAccount }}}&vadmin.Settings{Perms: user.GetPermissions()}{{{ else }}}&vadmin.Settings{}{{{ end }}}, ps, "admin")
		}
		ps.DefaultNavIcon = "cog"
		switch path[0] {
		case "server":
			info := util.DebugGetInfo(as.BuildInfo.Version, as.Started)
			ps.Data = info
			return controller.Render(w, r, as, &vadmin.ServerInfo{Info: info}, ps, "admin", "App Information"){{{ if .DangerousOK }}}
		case "cpu":
			switch path[1] {
			case "start":
				err := util.DebugStartCPUProfile()
				if err != nil {
					return "", err
				}
				return controller.FlashAndRedir(true, "started CPU profile", "/admin", w, ps)
			case "stop":
				pprof.StopCPUProfile()
				return controller.FlashAndRedir(true, "stopped CPU profile", "/admin", w, ps)
			default:
				return "", errors.Errorf("unhandled CPU profile action [%s]", path[1])
			}{{{ end }}}
		case "gc":
			timer := util.TimerStart()
			runtime.GC()
			msg := fmt.Sprintf("ran garbage collection in [%s]", timer.EndString())
			return controller.FlashAndRedir(true, msg, "/admin", w, ps){{{ if .DangerousOK }}}
		case "heap":
			err := util.DebugTakeHeapProfile()
			if err != nil {
				return "", err
			}
			return controller.FlashAndRedir(true, "wrote heap profile", "/admin", w, ps){{{ end }}}
		case "logs":
			ps.SetTitleAndData("Recent Logs", log.RecentLogs)
			return controller.Render(w, r, as, &vadmin.Logs{Logs: log.RecentLogs}, ps, "admin", "Recent Logs**folder")
		case "memusage":
			x := util.DebugMemStats()
			ps.Data = x
			return controller.Render(w, r, as, &vadmin.MemUsage{Mem: x}, ps, "admin", "Memory Usage**desktop"){{{ if .HasModule "migration" }}}
		case "migrations":
			ms := migrate.GetMigrations()
			am := migrate.ListMigrations(ps.Context, as.DB, nil, nil, ps.Logger)
			ps.Data = util.ValueMap{"available": ms, "applied": am}
			return controller.Render(w, r, as, &vadmin.Migrations{Available: ms, Applied: am}, ps, "admin", "Migrations**database"){{{ end }}}
		case "modules":
			di := util.DebugBuildInfo().Deps
			ps.SetTitleAndData("Modules", di)
			return controller.Render(w, r, as, &vadmin.Modules{Modules: di}, ps, "admin", "Modules**robot")
		case "request":
			ps.SetTitleAndData("Request Debug", cutil.RequestCtxToMap(w, r, as, ps))
			return controller.Render(w, r, as, &vadmin.Request{Req: r, Rsp: w}, ps, "admin", "Request**download")
		case "routes":
			ps.SetTitleAndData("HTTP Routes", cutil.AppRoutesList)
			return controller.Render(w, r, as, &vadmin.Routes{Routes: cutil.AppRoutesList}, ps, "admin", "Routes**folder")
		case "session":
			ps.SetTitleAndData("Session Debug", ps.Session)
			return controller.Render(w, r, as, &vadmin.Session{}, ps, "admin", "Session**play")
		case "sitemap":
			ps.SetTitleAndData("Sitemap", ps.Menu)
			return controller.Render(w, r, as, &vadmin.Sitemap{}, ps, "admin", "Sitemap**graph"){{{ if .HasModule "websocket" }}}
		case "sockets":
			return socketRoute(w, r, as, ps, path[1:]...){{{ end }}}
		// $PF_SECTION_START(admin-actions)$
		// $PF_SECTION_END(admin-actions)$
		default:
			return "", errors.Errorf("unhandled admin action [%s]", path[0])
		}
	})
}
