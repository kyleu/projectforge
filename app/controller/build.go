package controller

import (
	"github.com/kyleu/projectforge/views/vgit"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
)

func GitStatusAll(rc *fasthttp.RequestCtx) {
	act("build.status.all", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		statuses, err := as.Services.Git.GetStatusAll(prjs)
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve status")
		}

		ps.Data = statuses
		return render(rc, as, &vgit.StatusAll{Statuses: statuses}, ps, "projects", "Git")
	})
}
func GitStatus(rc *fasthttp.RequestCtx) {
	act("build.status", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := rcRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to load project")
		}
		status, err := as.Services.Git.GetStatus(prj)
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve status")
		}

		ps.Data = status
		return render(rc, as, &vgit.Status{Status: status}, ps, "projects", prj.Key, "Git")
	})
}
