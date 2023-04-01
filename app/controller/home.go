package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func Home(rc *fasthttp.RequestCtx) {
	Act("home", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		execs := as.Services.Exec.Execs
		mods := as.Services.Modules.Modules()
		ps.Data = util.ValueMap{"projects": prjs, "modules": mods}
		if string(rc.URI().QueryArgs().Peek("act")) == "clean" {
			err := Testbed(as, ps.Logger)
			if err != nil {
				return "", err
			}
		}
		return Render(rc, as, &views.Home{Projects: prjs, Execs: execs, Modules: mods}, ps)
	})
}

func Testbed(st *app.State, logger util.Logger) error {
	return nil
}
