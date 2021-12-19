package controller

import (
	"sort"
	"strings"
	"sync"

	"github.com/kyleu/projectforge/app/git"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/views/vgit"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
)

func GitActionAll(rc *fasthttp.RequestCtx) {
	a, _ := rcRequiredString(rc, "act", false)
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
		sort.Slice(results, func(i int, j int) bool {
			return strings.ToLower(results[i].Project.Title()) < strings.ToLower(results[j].Project.Title())
		})
		ps.Data = results
		return render(rc, as, &vgit.Results{Results: results}, ps, "projects", "Git")
	})
}

func gitStatusAll(prjs project.Projects, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (git.Results, error) {
	results := make(git.Results, 0, len(prjs))
	for _, prj := range prjs {
		s, err := as.Services.Git.Status(prj)
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
	var errs []error
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(prjs))
	for _, prj := range prjs {
		p := prj
		go func() {
			out, err := as.Services.Git.Fetch(p)
			mu.Lock()
			if err == nil {
				results = append(results, out)
			} else {
				errs = append(errs, errors.Wrapf(err, "can't fetch project [%s]", p.Key))
			}
			mu.Unlock()
			wg.Add(-1)
		}()
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return results, nil
}
