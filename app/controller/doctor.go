package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/views"
	"github.com/valyala/fasthttp"
)

func Doctor(rc *fasthttp.RequestCtx) {
	act("doctor", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		res := doctor.Check()
		ps.Data = res
		return render(rc, as, &views.Debug{}, ps, "doctor")
	})
}
