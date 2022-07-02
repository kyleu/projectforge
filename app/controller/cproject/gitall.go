package cproject

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/git"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/verror"
	"projectforge.dev/projectforge/views/vgit"
)

func GitActionAll(rc *fasthttp.RequestCtx) {
	a, _ := cutil.RCRequiredString(rc, "act", false)
	controller.Act("git.all."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) > 0 {
			prjs = prjs.WithTags(tags...)
		}

		var results git.Results
		var err error
		action := git.ActionStatusFromString(a)
		switch a {
		case git.ActionStatus.Key, "":
			results, err = gitStatusAll(prjs, rc, as, ps)
		case git.ActionFetch.Key:
			results, err = gitFetchAll(prjs, rc, as, ps)
		case git.ActionMagic.Key:
			argRes := cutil.CollectArgs(rc, gitMagicArgs)
			if len(argRes.Missing) > 0 {
				url := "/git/all/magic"
				ps.Data = argRes
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &verror.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes, Hidden: hidden}
				return controller.Render(rc, as, page, ps, "projects", "Git")
			}
			results, err = gitMagicAll(prjs, rc, as, ps)
		default:
			err = errors.Errorf("unhandled action [%s] for all projects", a)
		}
		if err != nil {
			return "", err
		}
		slices.SortFunc(results, func(l *git.Result, r *git.Result) bool {
			return strings.ToLower(l.Project.Title()) < strings.ToLower(r.Project.Title())
		})
		ps.Title = "[git] All Projects"
		ps.Data = results
		return controller.Render(rc, as, &vgit.Results{Action: action, Results: results, Projects: prjs, Tags: tags}, ps, "projects", "Git")
	})
}

func gitStatusAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Status(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitFetchAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(item *project.Project) (*git.Result, error) {
		return as.Services.Git.Fetch(ps.Context, item, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitMagicAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	message := string(rc.URI().QueryArgs().Peek("message"))
	dryRun := string(rc.URI().QueryArgs().Peek("dryRun")) == "true"
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Magic(ps.Context, prj, message, dryRun, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}
