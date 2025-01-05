package cproject

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vbuild"
)

func runCodeStats(prj *project.Project, res *action.Result, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	if res.HasErrors() {
		return "", errors.New(strings.Join(res.Errors, ", "))
	}
	stats, ok := res.Data.(*action.CodeStats)
	if !ok {
		return "", errors.Errorf("data is of type [%T], expected [Dependencies]", res.Data)
	}
	ps.SetTitleAndData(fmt.Sprintf("[%s] Code Stats", prj.Key), stats)
	page := &vbuild.CodeStats{Project: prj, Result: res, Stats: stats}
	return controller.Render(r, as, page, ps, "projects", prj.Key, "Dependency Management")
}

func runAllCodeStats(cfg util.ValueMap, prjs project.Projects, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	ret := map[string]*action.Result{}
	statsMap := map[string]*action.CodeStats{}

	for _, prj := range prjs {
		res := action.Apply(ps.Context, actionParams(prj.Key, action.TypeBuild, cfg, as, ps.Logger))
		if res.HasErrors() {
			return "", errors.New(strings.Join(res.Errors, ", "))
		}
		stats, ok := res.Data.(*action.CodeStats)
		if !ok {
			return "", errors.Errorf("data is of type [%T], expected [Pkgs]", res.Data)
		}
		ret[prj.Key] = res
		statsMap[prj.Key] = stats
	}

	ps.SetTitleAndData("Code Stats", statsMap)
	page := &vbuild.CodeStatsAll{Projects: prjs, Results: ret, CodeStats: statsMap}
	return controller.Render(r, as, page, ps, "projects", "Packages")
}
