package cproject

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vbuild"
)

func runPkgs(prj *project.Project, res *action.Result, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	if res.HasErrors() {
		return "", errors.New(strings.Join(res.Errors, ", "))
	}
	pkgs, ok := res.Data.(build.Pkgs)
	if !ok {
		return "", errors.Errorf("data is of type [%T], expected [Pkgs]", res.Data)
	}
	ps.SetTitleAndData(fmt.Sprintf("[%s] Packages", prj.Key), pkgs)
	return controller.Render(r, as, &vbuild.Packages{Project: prj, BuildResult: res, Packages: pkgs}, ps, "projects", prj.Key, "Packages")
}

func runAllPkgs(cfg util.ValueMap, prjs project.Projects, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	mu := sync.Mutex{}
	ret := map[string]*action.Result{}
	pkgs := map[string]build.Pkgs{}

	util.AsyncCollect(prjs, func(prj *project.Project) (string, error) {
		res := action.Apply(ps.Context, actionParams(prj.Key, action.TypeBuild, cfg, as, ps.Logger))
		packages, ok := res.Data.(build.Pkgs)
		if !ok {
			return "", errors.Errorf("data is of type [%T], expected [Pkgs]", res.Data)
		}
		mu.Lock()
		ret[prj.Key] = res
		pkgs[prj.Key] = packages
		mu.Unlock()
		return "", nil
	})

	ps.SetTitleAndData("Packages", pkgs)
	page := &vbuild.PackagesAll{Projects: prjs, Results: ret, Packages: pkgs}
	return controller.Render(r, as, page, ps, "projects", "Packages")
}
