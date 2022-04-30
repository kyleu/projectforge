package controller

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/git"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vgit"
)

func GitActionAll(rc *fasthttp.RequestCtx) {
	a, _ := RCRequiredString(rc, "act", false)
	act("git.all."+a, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		var results git.Results
		var err error
		switch a {
		case git.ActionStatus.Key, "":
			results, err = gitStatusAll(prjs, rc, as, ps)
		case git.ActionMagic.Key:
			results, err = gitMagicAll(prjs, rc, as, ps)
		case git.ActionFetch.Key:
			results, err = gitFetchAll(prjs, rc, as, ps)
		default:
			err = errors.Errorf("unhandled action [%s] for all projects", a)
		}
		if err != nil {
			return "", err
		}
		slices.SortFunc(results, func(l *git.Result, r *git.Result) bool {
			return strings.ToLower(l.Project.Title()) < strings.ToLower(r.Project.Title())
		})
		ps.Data = results
		return render(rc, as, &vgit.Results{Results: results}, ps, "projects", "Git")
	})
}

func gitStatusAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results := make(git.Results, 0, len(prjs))
	for _, prj := range prjs {
		s, err := as.Services.Git.Status(ps.Context, prj, ps.Logger)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get status for project [%s]", prj.Key)
		}
		results = append(results, s)
	}
	return results, nil
}

func gitMagicAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results := make(git.Results, 0, len(prjs))
	for _, prj := range prjs {
		out, err := as.Services.Git.Magic(prj)
		if err != nil {
			return nil, errors.Wrapf(err, "can't perform magic on project [%s]", prj.Key)
		}
		results = append(results, out)
	}
	return results, nil
}

func gitFetchAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results := make(git.Results, 0, len(prjs))

	results, errs := util.AsyncCollect(prjs, func(item *project.Project) (*git.Result, error) {
		return as.Services.Git.Fetch(ps.Context, item, ps.Logger)
	})
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return results, nil
}
