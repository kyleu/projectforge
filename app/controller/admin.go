// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/lib/user"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vadmin"
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
			mods, ok := debug.ReadBuildInfo()
			if !ok {
				return "", errors.New("unable to gather modules")
			}
			ps.Title = "Modules"
			ps.Data = mods.Deps
			return render(rc, as, &vadmin.Modules{Mods: mods.Deps}, ps, "admin", "Modules")
		case "request":
			ps.Title = "Request Debug"
			ps.Data = cutil.RequestCtxToMap(rc, nil)
			return render(rc, as, &vadmin.Request{RC: rc}, ps, "admin", "Request")
		case "session":
			ps.Title = "Session Debug"
			ps.Data = ps.Session
			return render(rc, as, &vadmin.Session{}, ps, "admin", "Session")
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
			return flashAndRedir(true, "wrote heap profile", "/admin", rc, ps)
		case "gc":
			start := util.TimerStart()
			runtime.GC()
			msg := fmt.Sprintf("ran garbage collection in [%s]", util.MicrosToMillis(util.TimerEnd(start)))
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
