package cproject

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vpage"
	"strings"

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
		if string(rc.URI().QueryArgs().Peek("loaded")) != "true" {
			page := &vpage.Load{URL: "/doctor/all?loaded=true", Title: "Running Doctor Checks", Message: "Hang tight..."}
			return controller.Render(rc, as, page, ps, "Welcome")
		}
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

func DoctorSolve(rc *fasthttp.RequestCtx) {
	controller.Act("doctor.solve", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		returnUrl := string(rc.URI().QueryArgs().Peek("return"))
		if returnUrl == "" {
			returnUrl = "/doctor/all"
		}

		c := checks.GetCheck(key)
		checks.CurrentModuleDeps = as.Services.Modules.Deps()
		ret := c.Check(ps.Context, ps.Logger)
		if len(ret.Solutions) == 0 {
			ps.Title = fmt.Sprintf("No solution available for [%s]", c.Title)
			return controller.Render(rc, as, &views.Debug{}, ps, "doctor", c.Title)
		}
		execs := lo.FilterMap(ret.Solutions, func(sol string, _ int) (string, bool) {
			if !strings.HasPrefix(sol, "!") {
				return "", false
			}
			return strings.TrimPrefix(sol, "!"), true
		})
		if len(execs) == 0 {
			ps.Title = fmt.Sprintf("No solution for [%s] can be solved automatically", c.Title)
			return controller.Render(rc, as, &views.Debug{}, ps, "doctor", c.Title)
		}
		for idx, ex := range execs {
			exec := as.Services.Exec.NewExec(fmt.Sprintf("solve-%s-%d", c.Key, idx), ex, ".")
			err = exec.Start()
			if err != nil {
				return "", errors.Wrapf(err, "unable to run [%s]", ex)
			}
			err = exec.Wait()
			if err != nil {
				return "", errors.Wrapf(err, "error running [%s]", ex)
			}
		}

		return controller.FlashAndRedir(true, "Issue solved", returnUrl, rc, ps)
	})
}
