package controller

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/build"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vbuild"
)

func actionParams(tgt string, t action.Type, cfg util.ValueMap, as *app.State, logger util.Logger) *action.Params {
	return &action.Params{
		ProjectKey: tgt, T: t, Cfg: cfg,
		MSvc: as.Services.Modules, PSvc: as.Services.Projects, ESvc: as.Services.Export, Logger: logger,
	}
}

func runDeps(prj *project.Project, res *action.Result, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	if res.HasErrors() {
		return "", errors.Errorf(strings.Join(res.Errors, ", "))
	}
	deps, ok := res.Data.(build.Dependencies)
	if !ok {
		return "", errors.Errorf("data is of type [%T], expected [Dependencies]", res.Data)
	}
	return render(rc, as, &vbuild.Deps{Project: prj, BuildResult: res, Dependencies: deps}, ps, "projects", prj.Key, "Dependency Management")
}

func runAllDeps(cfg util.ValueMap, prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	var msg string
	key := cfg.GetStringOpt("key")
	if pj := cfg.GetStringOpt("project"); pj != "" {
		result, err := build.SetDepsProject(ps.Context, prjs, pj, as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}
		msg = result
	} else {
		version := cfg.GetStringOpt("version")
		if key != "" && version != "" {
			result, err := build.SetDepsMap(ps.Context, prjs, &build.Dependency{Key: key, Version: version}, as.Services.Projects, ps.Logger)
			if err != nil {
				return "", err
			}
			msg = result
		}
	}

	ret, err := build.LoadDepsMap(prjs, 2, as.Services.Projects)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	ps.Data = ret
	page := &vbuild.DepMap{Message: msg, Result: ret}
	return render(rc, as, page, ps, "projects", "Dependencies")
}