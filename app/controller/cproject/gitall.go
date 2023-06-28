package cproject

import (
	"strconv"
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
	"projectforge.dev/projectforge/views/vgit"
	"projectforge.dev/projectforge/views/vpage"
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
			results, err = gitStatusAll(prjs, as, ps)
		case git.ActionFetch.Key:
			results, err = gitFetchAll(prjs, as, ps)
		case git.ActionPull.Key:
			results, err = gitPullAll(prjs, as, ps)
		case git.ActionOutdated.Key:
			results, err = gitOutdatedAll(prjs, as, ps)
		case git.ActionHistory.Key:
			argRes := cutil.CollectArgs(rc, gitHistoryArgs)
			if argRes.HasMissing() {
				url := "/git/all/history"
				ps.Data = argRes
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &vpage.Args{URL: url, Directions: "Choose your options", ArgRes: argRes, Hidden: hidden}
				return controller.Render(rc, as, page, ps, "projects", "Git")
			}
			results, err = gitHistoryAll(prjs, rc, as, ps)
		case git.ActionMagic.Key:
			argRes := cutil.CollectArgs(rc, gitMagicArgs)
			if argRes.HasMissing() {
				url := "/git/all/magic"
				ps.Data = argRes
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &vpage.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes, Hidden: hidden}
				return controller.Render(rc, as, page, ps, "projects", "Git")
			}
			results, err = gitMagicAll(prjs, rc, as, ps)
		case git.ActionUndoCommit.Key:
			results, err = gitUndoAll(prjs, as, ps)
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

func gitStatusAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Status(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitFetchAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Fetch(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitPullAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Pull(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitOutdatedAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Outdated(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitHistoryAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	path := string(rc.URI().QueryArgs().Peek("path"))
	since, _ := util.TimeFromString(string(rc.URI().QueryArgs().Peek("since")))
	authors := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("authors")), ",")
	limit, _ := strconv.Atoi(string(rc.URI().QueryArgs().Peek("limit")))
	commit := string(rc.URI().QueryArgs().Peek("commit"))
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		hist := &git.HistoryResult{Path: path, Since: since, Authors: authors, Commit: commit, Limit: limit}
		return as.Services.Git.History(ps.Context, prj, hist, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitMagicAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	message := string(rc.URI().QueryArgs().Peek("message"))
	dryRun := cutil.QueryStringBool(rc, "dryRun")
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Magic(ps.Context, prj, message, dryRun, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitUndoAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.UndoCommit(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}
