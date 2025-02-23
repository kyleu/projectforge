package clib

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vadmin"
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
			return controller.Render(r, as, &vadmin.Settings{BuildInfo: as.BuildInfo}, ps, "admin")
		}
		ps.DefaultNavIcon = "cog"
		switch path[0] {
		case "server":
			info := util.DebugGetInfo(as.BuildInfo.Version, as.Started)
			ps.SetTitleAndData("Server Info", info)
			return controller.Render(r, as, &vadmin.ServerInfo{Info: info}, ps, "admin", "App Information")
		case "cpu":
			switch path[1] {
			case util.KeyStart:
				err := util.DebugStartCPUProfile()
				if err != nil {
					return "", err
				}
				return controller.FlashAndRedir(true, "started CPU profile", "/admin", ps)
			case "stop":
				pprof.StopCPUProfile()
				return controller.FlashAndRedir(true, "stopped CPU profile", "/admin", ps)
			default:
				return "", errors.Errorf("unhandled CPU profile action [%s]", path[1])
			}
		case "gc":
			timer := util.TimerStart()
			runtime.GC()
			msg := fmt.Sprintf("ran garbage collection in [%s]", timer.EndString())
			return controller.FlashAndRedir(true, msg, "/admin", ps)
		case "heap":
			err := util.DebugTakeHeapProfile()
			if err != nil {
				return "", err
			}
			return controller.FlashAndRedir(true, "wrote heap profile", "/admin", ps)
		case "logs":
			ps.SetTitleAndData("Recent Logs", log.RecentLogs)
			return controller.Render(r, as, &vadmin.Logs{Logs: log.RecentLogs}, ps, "admin", "Recent Logs**folder")
		case "memusage":
			x := util.DebugMemStats()
			ps.SetTitleAndData("Memory Usage", x)
			return controller.Render(r, as, &vadmin.MemUsage{Mem: x}, ps, "admin", "Memory Usage**desktop")
		case "modules":
			di := util.DebugBuildInfo().Deps
			ps.SetTitleAndData("Modules", di)
			return controller.Render(r, as, &vadmin.Modules{Modules: di, Version: as.BuildInfo.Version}, ps, "admin", "Modules**robot")
		case "request":
			ps.SetTitleAndData("Request Debug", cutil.RequestCtxToMap(r, as, ps))
			return controller.Render(r, as, &vadmin.Request{Req: r, Rsp: w}, ps, "admin", "Request**download")
		case "routes":
			ps.SetTitleAndData("HTTP Routes", cutil.AppRoutesList)
			return controller.Render(r, as, &vadmin.Routes{Routes: cutil.AppRoutesList}, ps, "admin", "Routes**folder")
		case "session":
			ps.SetTitleAndData("Session Debug", ps.Session)
			return controller.Render(r, as, &vadmin.Session{}, ps, "admin", "Session**play")
		case "sitemap":
			ps.SetTitleAndData("Sitemap", ps.Menu)
			return controller.Render(r, as, &vadmin.Sitemap{}, ps, "admin", "Sitemap**graph")
		case "sockets":
			return socketRoute(w, r, as, ps, path[1:]...)
		// $PF_SECTION_START(admin-actions)$
		// $PF_SECTION_END(admin-actions)$
		default:
			return "", errors.Errorf("unhandled admin action [%s]", path[0])
		}
	})
}
