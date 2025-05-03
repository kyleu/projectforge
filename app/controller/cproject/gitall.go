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
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/project"
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
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Status(ps.Context, ps.Logger)
			})
		case git.ActionFetch.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Fetch(ps.Context, ps.Logger)
			})
		case git.ActionPull.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Pull(ps.Context, ps.Logger)
			})
		case git.ActionPush.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Push(ps.Context, ps.Logger)
			})
		case git.ActionReset.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Reset(ps.Context, ps.Logger)
			})
		case git.ActionOutdated.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).Outdated(ps.Context, ps.Logger)
			})
		case git.ActionHistory.Key:
			argRes := util.FieldDescsCollect(r, gitHistoryArgs)
			if argRes.HasMissing() {
				url := "/git/all/history"
				ps.SetTitleAndData("Git History", argRes)
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				page := &vpage.Args{URL: url, Directions: "Choose your options", Results: argRes, Hidden: hidden}
				return controller.Render(r, as, page, ps, "projects", "Git**git")
			}
			results, err = gitHistoryAll(prjs, r, ps)
		case git.ActionMagic.Key:
			argRes := util.FieldDescsCollect(r, gitMagicArgs)
			if argRes.HasMissing() {
				url := "/git/all/magic"
				ps.SetTitleAndData("Git Magic!", argRes)
				hidden := map[string]string{"tags": strings.Join(tags, ",")}
				warning := "Are you sure you'd like to commit and push for all projects?"
				page := &vpage.Args{URL: url, Directions: "Enter your commit message", Results: argRes, Hidden: hidden, Warning: warning}
				return controller.Render(r, as, page, ps, "projects", "Git**git")
			}
			results, err = gitMagicAll(prjs, r, ps)
		case git.ActionUndoCommit.Key:
			results, err = gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
				return git.NewService(prj.Key, prj.Path).UndoCommit(ps.Context, ps.Logger)
			})
		default:
			err = errors.Errorf("unhandled action [%s] for all projects", a)
		}
		if err != nil {
			return "", err
		}
		slices.SortFunc(results, func(l *git.Result, r *git.Result) int {
			lp, rp := prjs.Get(l.Project), prjs.Get(r.Project)
			return cmp.Compare(strings.ToLower(lp.Title()), strings.ToLower(rp.Title()))
		})
		ps.SetTitleAndData("[git] All Projects", results)
		return controller.Render(r, as, &vgit.Results{Action: action, Results: results, Projects: prjs, Tags: tags}, ps, "projects", "Git")
	})
}

func gitAll(prjs project.Projects, f func(prj *project.Project) (*git.Result, error)) (git.Results, error) {
	results, errs := util.AsyncCollect(prjs, f)
	return results, util.ErrorMerge(errs...)
}

func gitHistoryAll(prjs project.Projects, r *http.Request, ps *cutil.PageState) (git.Results, error) {
	path := r.URL.Query().Get("path")
	since, _ := util.TimeFromString(r.URL.Query().Get("since"))
	authors := util.StringSplitAndTrim(r.URL.Query().Get("authors"), ",")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	commit := r.URL.Query().Get("commit")
	return gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
		args := &git.HistoryArgs{Path: path, Since: since, Authors: authors, Commit: commit, Limit: int(limit)}
		return git.NewService(prj.Key, prj.Path).History(ps.Context, args, ps.Logger)
	})
}

func gitMagicAll(prjs project.Projects, r *http.Request, ps *cutil.PageState) (git.Results, error) {
	message := r.URL.Query().Get("message")
	dryRun := cutil.QueryStringBool(r, "dryRun")
	return gitAll(prjs, func(prj *project.Project) (*git.Result, error) {
		return git.NewService(prj.Key, prj.Path).Magic(ps.Context, message, dryRun, ps.Logger)
	})
}
