package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vbuild"
)

func runCoverage(prj *project.Project, res *action.Result, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	if res.HasErrors() {
		return "", errors.New(util.StringJoin(res.Errors, ", "))
	}
	coverage, err := util.Cast[*action.Coverage](res.Data)
	if err != nil {
		return "", err
	}
	ps.SetTitleAndData(fmt.Sprintf("[%s] Code Coverage", prj.Key), coverage)
	page := &vbuild.Coverage{Project: prj, Result: res, Coverage: coverage}
	return controller.Render(r, as, page, ps, "projects", prj.Key, "Code Coverage")
}

func runAllCoverage(cfg util.ValueMap, prjs project.Projects, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	ret := map[string]*action.Result{}
	coverageMap := map[string]*action.Coverage{}

	for _, prj := range prjs {
		res := action.Apply(ps.Context, actionParams(prj.Key, action.TypeBuild, cfg, as, ps.Logger))
		if res.HasErrors() {
			return "", errors.New(util.StringJoin(res.Errors, ", "))
		}
		coverage, err := util.Cast[*action.Coverage](res.Data)
		if err != nil {
			return "", err
		}
		ret[prj.Key] = res
		coverageMap[prj.Key] = coverage
	}

	ps.SetTitleAndData("Code Coverage", coverageMap)
	page := &vbuild.CoverageAll{Projects: prjs, Results: ret, Coverage: coverageMap}
	return controller.Render(r, as, page, ps, "projects", "Packages")
}
