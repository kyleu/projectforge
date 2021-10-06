package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/doctor/checks"
	"github.com/kyleu/projectforge/views/vdoctor"
	"github.com/valyala/fasthttp"
)

func Doctor(rc *fasthttp.RequestCtx) {
	act("doctor", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ret := checks.ForModules(prjs.AllModules())
		ps.Data = ret
		return render(rc, as, &vdoctor.List{Checks: ret}, ps, "doctor")
	})
}

func DoctorRunAll(rc *fasthttp.RequestCtx) {
	act("doctor.run.all", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ret := checks.CheckAll(prjs.AllModules())
		ps.Data = ret
		return render(rc, as, &vdoctor.Results{Results: ret}, ps, "doctor", "All")
	})
}

func DoctorRun(rc *fasthttp.RequestCtx) {
	act("doctor.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := rcRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		c := checks.GetCheck(key)
		ret := c.Check()
		ps.Data = ret
		return render(rc, as, &vdoctor.Results{Results: doctor.Results{ret}}, ps, "doctor", c.Title)
	})
}
