package ctest

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vaction"
)

func bootstrapTest(as *app.State, ps *cutil.PageState) (layout.Page, error) {
	cfg := util.ValueMap{}
	cfg.Add("path", "./testproject", "method", "bootstrap", "wipe", true)
	res := action.Apply(ps.Context, bootstrapParams("testproject", action.TypeTest, cfg, as, ps.Logger))

	ps.SetTitleAndData("Bootstrap", res)

	if _, err := as.Services.Projects.Refresh(ps.Logger); err != nil {
		return nil, err
	}

	prj, err := as.Services.Projects.Get("testproject")
	if err != nil {
		return nil, err
	}

	page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: res}}
	return page, nil
}

func bootstrapParams(tgt string, t action.Type, cfg util.ValueMap, as *app.State, logger util.Logger) *action.Params {
	return &action.Params{
		ProjectKey: tgt, T: t, Cfg: cfg,
		MSvc: as.Services.Modules, PSvc: as.Services.Projects, XSvc: as.Services.Exec, SSvc: as.Services.Socket, ESvc: as.Services.Export, Logger: logger,
	}
}
