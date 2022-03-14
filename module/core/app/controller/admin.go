package controller

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database/migrate"{{{ end }}}
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vadmin"
)

func Admin(rc *fasthttp.RequestCtx) {
	act("admin", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		path := util.StringSplitAndTrim(strings.TrimPrefix(string(rc.URI().Path()), "/admin"), "/")
		if len(path) == 0 {
			ps.Title = "Administration"
			x := &runtime.MemStats{}
			runtime.ReadMemStats(x)
			return render(rc, as, &vadmin.Settings{Perms: user.GetPermissions(), Mem: x}, ps, "admin")
		}
		switch path[0] {
		case "modules":
			di, err := util.GetDebugInfo()
			if err != nil {
				return "", err
			}
			ps.Title = "Modules"
			ps.Data = di
			return render(rc, as, &vadmin.Modules{Info: di}, ps, "admin", "Modules")
		case "request":
			ps.Title = "Request Debug"
			ps.Data = cutil.RequestCtxToMap(rc, nil)
			return render(rc, as, &vadmin.Request{RC: rc}, ps, "admin", "Request")
		case "session":
			ps.Title = "Session Debug"
			ps.Data = ps.Session
			return render(rc, as, &vadmin.Session{}, ps, "admin", "Session"){{{ if .HasModule "queue" }}}
		case "queue":
			ps.Title = "Queue Debug"
			ps.Data = as.Services.Publisher
			return render(rc, as, &vadmin.Queue{Con: as.Services.Consumer, Pub: as.Services.Publisher}, ps, "admin", "Queue"){{{ end }}}
		case "cpu":
			switch path[1] {
			case "start":
				err := startCPUProfile()
				if err != nil {
					return "", err
				}
				return flashAndRedir(true, "started CPU profile", "/admin", rc, ps)
			case "stop":
				pprof.StopCPUProfile()
				return flashAndRedir(true, "stopped CPU profile", "/admin", rc, ps)
			default:
				return "", errors.Errorf("unhandled CPU profile action [%s]", path[1])
			}
		case "heap":
			err := takeHeapProfile()
			if err != nil {
				return "", err
			}
			return flashAndRedir(true, "wrote heap profile", "/admin", rc, ps){{{ if .HasModule "migration" }}}
		case "migrations":
			ms := migrate.GetMigrations()
			am := migrate.ListMigrations(ps.Context, as.DB, nil, as.Logger)
			ps.Data = map[string]interface{}{"available": ms, "applied": am}
			return render(rc, as, &vadmin.Migrations{Available: ms, Applied: am}, ps, "admin", "Migrations"){{{ end }}}
		case "gc":
			timer := util.TimerStart()
			runtime.GC()
			msg := fmt.Sprintf("ran garbage collection in [%s]", timer.EndString())
			return flashAndRedir(true, msg, "/admin", rc, ps)
		// $PF_SECTION_START(admin-actions)$
		// $PF_SECTION_END(admin-actions)$

		default:
			return "", errors.Errorf("unhandled admin action [%s]", path[0])
		}
	})
}

func startCPUProfile() error {
	f, err := os.Create("./tmp/cpu.pprof")
	if err != nil {
		return err
	}
	return pprof.StartCPUProfile(f)
}

func takeHeapProfile() error {
	f, err := os.Create("./tmp/mem.pprof")
	if err != nil {
		return err
	}
	return pprof.WriteHeapProfile(f)
}
