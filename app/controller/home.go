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
		if x := string(rc.URI().QueryArgs().Peek("test")); x != "" {
			err := Testbed(as, x, ps.Logger)
			if err != nil {
				return "", err
			}
		}
		return Render(rc, as, &views.Home{Projects: prjs, Execs: execs, Modules: mods}, ps)
	})
}

func Testbed(st *app.State, prjKey string, logger util.Logger) error {
	svcs := st.Services
	prj, err := svcs.Projects.Get(prjKey)
	if err != nil {
		return err
	}
	args, err := prj.ModuleArgExport(svcs.Projects, logger)
	if err != nil {
		return err
	}
	for _, m := range args.Models {
		m.AddTag("search")
		for _, col := range m.Columns {
			if col.Name == "DateCreated" {
				col.AddTag("created")
			}
			if col.PK {
				col.Search = true
				col.RemoveTag("search")
			}
		}
		switch m.Name {
		case "Bar":
			// m.Group = []string{"foo"}
		}
	}

	err = svcs.Projects.Save(prj, logger)
	if err != nil {
		return err
	}

	return nil
}
