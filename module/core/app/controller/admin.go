package controller

import (
	"os"
	"runtime/debug"
	"runtime/pprof"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/database/migrate"{{{ end }}}
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vadmin"
)

func Admin(rc *fasthttp.RequestCtx) {
	act("admin", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		path := util.SplitAndTrim(strings.TrimPrefix(string(rc.URI().Path()), "/admin"), "/")
		if len(path) == 0 {
			ps.Title = "Administration"
			return render(rc, as, &vadmin.List{}, ps, "Administration")
		}
		switch path[0] {
		case "modules":
			mods, ok := debug.ReadBuildInfo()
			if !ok {
				return "", errors.New("unable to gather modules")
			}
			ps.Title = "Modules"
			ps.Data = mods.Deps
			return render(rc, as, &vadmin.Modules{Mods: mods.Deps}, ps, "Administration||/admin", "Modules")
		case "session":
			ps.Title = "Session Debug"
			ps.Data = ps.Session
			return render(rc, as, &vadmin.Session{}, ps, "Administration||/admin", "Session")
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
			ps.Data = ms
			return render(rc, as, &vadmin.Migrations{Migrations: ms}, ps, "Administration||/admin", "Migrations"){{{ end }}}

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
