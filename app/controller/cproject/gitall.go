package cproject

import (
	"cmp"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/git"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vgit"
	"projectforge.dev/projectforge/views/vpage"
)

func GitActionAll(w http.ResponseWriter, r *http.Request) {
	a, _ := cutil.PathString(r, "act", false)
	controller.Act("git.all."+a, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(r.URL.Query().Get("tags"), ",")
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
			argRes := cutil.CollectArgs(r, gitHistoryArgs)
			if argRes.HasMissing() {
				url := "/git/all/history"
				ps.Data = argRes
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &vpage.Args{URL: url, Directions: "Choose your options", ArgRes: argRes, Hidden: hidden}
				return controller.Render(w, r, as, page, ps, "projects", "Git")
			}
			results, err = gitHistoryAll(prjs, r, as, ps)
		case git.ActionMagic.Key:
			argRes := cutil.CollectArgs(r, gitMagicArgs)
			if argRes.HasMissing() {
				url := "/git/all/magic"
				ps.Data = argRes
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &vpage.Args{URL: url, Directions: "Enter your commit message", ArgRes: argRes, Hidden: hidden}
				return controller.Render(w, r, as, page, ps, "projects", "Git")
			}
			results, err = gitMagicAll(prjs, r, as, ps)
		case git.ActionPush.Key:
			results, err = gitPushAll(prjs, as, ps)
		case git.ActionUndoCommit.Key:
			results, err = gitUndoAll(prjs, as, ps)
		default:
			err = errors.Errorf("unhandled action [%s] for all projects", a)
		}
		if err != nil {
			return "", err
		}
		slices.SortFunc(results, func(l *git.Result, r *git.Result) int {
			return cmp.Compare(strings.ToLower(l.Project.Title()), strings.ToLower(r.Project.Title()))
		})
		ps.SetTitleAndData("[git] All Projects", results)
		return controller.Render(w, r, as, &vgit.Results{Action: action, Results: results, Projects: prjs, Tags: tags}, ps, "projects", "Git")
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

func gitPushAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Push(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitOutdatedAll(prjs project.Projects, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		return as.Services.Git.Outdated(ps.Context, prj, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitHistoryAll(prjs project.Projects, r *http.Request, as *app.State, ps *cutil.PageState) (git.Results, error) {
	path := r.URL.Query().Get("path")
	since, _ := util.TimeFromString(r.URL.Query().Get("since"))
	authors := util.StringSplitAndTrim(r.URL.Query().Get("authors"), ",")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	commit := r.URL.Query().Get("commit")
	results, errs := util.AsyncCollect(prjs, func(prj *project.Project) (*git.Result, error) {
		hist := &git.HistoryResult{Path: path, Since: since, Authors: authors, Commit: commit, Limit: int(limit)}
		return as.Services.Git.History(ps.Context, prj, hist, ps.Logger)
	})
	return results, util.ErrorMerge(errs...)
}

func gitMagicAll(prjs project.Projects, r *http.Request, as *app.State, ps *cutil.PageState) (git.Results, error) {
	message := r.URL.Query().Get("message")
	dryRun := cutil.QueryStringBool(r, "dryRun")
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
