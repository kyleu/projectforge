// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vadmin"
)

var AppRoutesList map[string][]string

func Admin(rc *fasthttp.RequestCtx) {
	path := util.StringSplitAndTrim(strings.TrimPrefix(string(rc.URI().Path()), "/admin"), "/")
	key := "admin"
	if len(path) > 0 {
		key += "." + strings.Join(path, ".")
	}
	controller.Act(key, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		if len(path) == 0 {
			ps.Title = "Administration"
			ps.Data = "administration"
			return controller.Render(rc, as, &vadmin.Settings{}, ps, "admin")
		}
		switch path[0] {
		case "server":
			info := util.DebugGetInfo(as.BuildInfo.Version, as.Started)
			ps.Data = info
			return controller.Render(rc, as, &vadmin.ServerInfo{Info: info}, ps, "admin", "App Information")
		case "cpu":
			switch path[1] {
			case "start":
				err := util.DebugStartCPUProfile()
				if err != nil {
					return "", err
				}
				return controller.FlashAndRedir(true, "started CPU profile", "/admin", rc, ps)
			case "stop":
				pprof.StopCPUProfile()
				return controller.FlashAndRedir(true, "stopped CPU profile", "/admin", rc, ps)
			default:
				return "", errors.Errorf("unhandled CPU profile action [%s]", path[1])
			}
		case "gc":
			timer := util.TimerStart()
			runtime.GC()
			msg := fmt.Sprintf("ran garbage collection in [%s]", timer.EndString())
			return controller.FlashAndRedir(true, msg, "/admin", rc, ps)
		case "heap":
			err := util.DebugTakeHeapProfile()
			if err != nil {
				return "", err
			}
			return controller.FlashAndRedir(true, "wrote heap profile", "/admin", rc, ps)
		case "logs":
			ps.Title = "Recent Logs"
			ps.Data = log.RecentLogs
			return controller.Render(rc, as, &vadmin.Logs{Logs: log.RecentLogs}, ps, "admin", "Recent Logs")
		case "memusage":
			x := util.DebugMemStats()
			ps.Data = x
			return controller.Render(rc, as, &vadmin.MemUsage{Mem: x}, ps, "admin", "Memory Usage")
		case "modules":
			di := util.DebugBuildInfo().Deps
			ps.Title = "Modules"
			ps.Data = di
			return controller.Render(rc, as, &vadmin.Modules{Modules: di}, ps, "admin", "Modules")
		case "request":
			ps.Title = "Request Debug"
			ps.Data = cutil.RequestCtxToMap(rc, as, ps)
			return controller.Render(rc, as, &vadmin.Request{RC: rc}, ps, "admin", "Request")
		case "routes":
			ps.Title = "HTTP Routes"
			ps.Data = AppRoutesList
			return controller.Render(rc, as, &vadmin.Routes{Routes: AppRoutesList}, ps, "admin", "Routes")
		case "session":
			ps.Title = "Session Debug"
			ps.Data = ps.Session
			return controller.Render(rc, as, &vadmin.Session{}, ps, "admin", "Session")
		case "sitemap":
			ps.Title = "Sitemap"
			ps.Data = ps.Menu
			return controller.Render(rc, as, &vadmin.Sitemap{}, ps, "admin", "Sitemap")
		case "sockets":
			return socketRoute(rc, as, ps, path[1:]...)
		// $PF_SECTION_START(admin-actions)$
		// $PF_SECTION_END(admin-actions)$
		default:
			return "", errors.Errorf("unhandled admin action [%s]", path[0])
		}
	})
}
