package cproject

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/build"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vbuild"
)

func runPkgs(prj *project.Project, res *action.Result, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	if res.HasErrors() {
		return "", errors.Errorf(strings.Join(res.Errors, ", "))
	}
	pkgs, ok := res.Data.(build.Pkgs)
	if !ok {
		return "", errors.Errorf("data is of type [%T], expected [Pkgs]", res.Data)
	}
	ps.Title = fmt.Sprintf("[%s] Dependencies", prj.Key)
	ps.Data = pkgs
	return controller.Render(rc, as, &vbuild.Packages{Project: prj, BuildResult: res, Packages: pkgs}, ps, "projects", prj.Key, "Packages")
}

func runAllPkgs(cfg util.ValueMap, prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
	if len(tags) > 0 {
		prjs = prjs.WithTags(tags...)
	}
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
	ps.Title = "Dependency Merge"
	ps.Data = ret
	page := &vbuild.DepMap{Message: msg, Result: ret, Tags: tags}
	return controller.Render(rc, as, page, ps, "projects", "Dependencies")
}
