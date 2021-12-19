package controller

import (
	"github.com/kyleu/projectforge/views/vbuild"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
)

func BuildList(rc *fasthttp.RequestCtx) {
	act("build.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		statuses, err := as.Services.Build.GetStatusAll(prjs)
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve status")
		}

		ps.Data = statuses
		return render(rc, as, &vbuild.StatusAll{Statuses: statuses}, ps, "projects", "Build")
	})
}
