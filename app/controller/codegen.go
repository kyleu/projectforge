package controller

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/views"
)

func Codegen(rc *fasthttp.RequestCtx) {
	act("codegen", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		ps.Data = prj
		return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Code Generation")
	})
}

func CodegenAct(rc *fasthttp.RequestCtx) {
	act("codegen.act", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		action, err := rcRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}

		switch action {
		case "test":
			ps.Data = prj
			return render(rc, as, &views.Debug{}, ps, "projects", prj.Key, "Code Generation")
		default:
			return "", errors.Errorf("invalid codegen action [%s]", action)
		}
	})
}
