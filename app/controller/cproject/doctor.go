package cproject

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/views/vdoctor"
)

func Doctor(rc *fasthttp.RequestCtx) {
	controller.Act("doctor", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ret := checks.ForModules(prjs.AllModules())
		ps.Title = "Doctor"
		ps.Data = ret
		return controller.Render(rc, as, &vdoctor.List{Checks: ret}, ps, "doctor")
	})
}

func DoctorRunAll(rc *fasthttp.RequestCtx) {
	controller.Act("doctor.run.all", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		checks.CurrentModuleDeps = as.Services.Modules.Deps()
		ret := checks.CheckAll(ps.Context, prjs.AllModules(), ps.Logger)
		ps.Title = "Doctor Results"
		ps.Data = ret
		return controller.Render(rc, as, &vdoctor.Results{Results: ret}, ps, "doctor", "All")
	})
}

func DoctorRun(rc *fasthttp.RequestCtx) {
	controller.Act("doctor.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		c := checks.GetCheck(key)
		checks.CurrentModuleDeps = as.Services.Modules.Deps()
		ret := c.Check(ps.Context, ps.Logger)
		ps.Title = c.Title + " Result"
		ps.Data = ret
		return controller.Render(rc, as, &vdoctor.Results{Results: doctor.Results{ret}}, ps, "doctor", c.Title)
	})
}
