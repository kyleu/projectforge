package controller

import (
	"github.com/kyleu/projectforge/app/git"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/views"
	"github.com/kyleu/projectforge/views/vgit"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
)

func GitAction(rc *fasthttp.RequestCtx) {
	a, _ := rcRequiredString(rc, "act", false)
	act("git.action."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := rcRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		switch a {
		case git.ActionStatus.Key, "":
			return gitStatus(prj, rc, as, ps)
		case git.ActionCreateRepo.Key:
			return gitCreateRepo(prj, rc, as, ps)
		case git.ActionMagic.Key:
			return gitMagic(prj, rc, as, ps)
		case git.ActionFetch.Key:
			return gitFetch(prj, rc, as, ps)
		case git.ActionCommit.Key:
			return gitCommit(prj, rc, as, ps)
		default:
			return "", errors.Errorf("unhandled action [%s]", a)
		}
	})
}

func gitStatus(prj *project.Project, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	result, err := as.Services.Git.Status(prj)
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve status")
	}
	ps.Data = result
	return render(rc, as, &vgit.Result{Result: result}, ps, "projects", prj.Key, "Git")
}

func gitCreateRepo(prj *project.Project, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	result, err := as.Services.Git.CreateRepo(prj)
	if err != nil {
		return "", errors.Wrap(err, "unable to create repo")
	}

	ps.Data = result
	return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Git")
}

func gitMagic(prj *project.Project, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	result, err := as.Services.Git.Magic(prj)
	if err != nil {
		return "", errors.Wrap(err, "unable to perform magic on repo")
	}

	ps.Data = result
	return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Git")
}

func gitFetch(prj *project.Project, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	out, err := as.Services.Git.Fetch(prj)
	if err != nil {
		return "", errors.Wrap(err, "unable to fetch repo")
	}

	ps.Data = out
	return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Git")
}

func gitCommit(prj *project.Project, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	out, err := as.Services.Git.Commit(prj, "TODO")
	if err != nil {
		return "", errors.Wrap(err, "unable to commit repo")
	}

	ps.Data = out
	return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Git")
}
