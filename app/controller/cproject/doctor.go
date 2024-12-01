package cproject

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vdoctor"
	"projectforge.dev/projectforge/views/vpage"
)

func Doctor(w http.ResponseWriter, r *http.Request) {
	controller.Act("doctor", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ret := checks.ForModules(prjs.AllModules())
		ps.SetTitleAndData("Doctor", ret)
		return controller.Render(r, as, &vdoctor.List{Checks: ret}, ps, "doctor")
	})
}

func DoctorRunAll(w http.ResponseWriter, r *http.Request) {
	controller.Act("doctor.run.all", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		if r.URL.Query().Get("loaded") != util.BoolTrue {
			page := &vpage.Load{URL: "/doctor/all?loaded=true", Title: "Running Doctor Checks", Message: "Hang tight..."}
			return controller.Render(r, as, page, ps, "Welcome")
		}
		prjs := as.Services.Projects.Projects()
		checks.SetModules(as.Services.Modules.Deps(), as.Services.Modules.Dangerous())
		ret := checks.CheckAll(ps.Context, prjs.AllModules(), ps.Logger)
		ps.SetTitleAndData("Doctor Results", ret)
		return controller.Render(r, as, &vdoctor.Results{Results: ret}, ps, "doctor", "All")
	})
}

func DoctorRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("doctor.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		c := checks.GetCheck(key)
		checks.SetModules(as.Services.Modules.Deps(), as.Services.Modules.Dangerous())
		ret := c.Check(ps.Context, ps.Logger)
		ps.SetTitleAndData(c.Title+" Result", ret)
		return controller.Render(r, as, &vdoctor.Results{Results: doctor.Results{ret}}, ps, "doctor", c.Title)
	})
}

func DoctorSolve(w http.ResponseWriter, r *http.Request) {
	controller.Act("doctor.solve", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		returnURL := util.OrDefault(r.URL.Query().Get("return"), "/doctor/all")
		c := checks.GetCheck(key)
		checks.SetModules(as.Services.Modules.Deps(), as.Services.Modules.Dangerous())
		ret := c.Check(ps.Context, ps.Logger)
		if len(ret.Solutions) == 0 {
			ps.SetTitleAndData(fmt.Sprintf("No solution available for [%s]", c.Title), c)
			return controller.Render(r, as, &views.Debug{}, ps, "doctor", c.Title)
		}
		execs := lo.FilterMap(ret.Solutions, func(sol string, _ int) (string, bool) {
			if !strings.HasPrefix(sol, "!") {
				return "", false
			}
			return strings.TrimPrefix(sol, "!"), true
		})
		if len(execs) == 0 {
			ps.SetTitleAndData(fmt.Sprintf("No solution for [%s] can be solved automatically", c.Title), c)
			return controller.Render(r, as, &views.Debug{}, ps, "doctor", c.Title)
		}
		for idx, ex := range execs {
			exec := as.Services.Exec.NewExec(fmt.Sprintf("solve-%s-%d", c.Key, idx), ex, ".", false)
			err = exec.Start()
			if err != nil {
				return "", errors.Wrapf(err, "unable to run [%s]", ex)
			}
			err = exec.Wait()
			if err != nil {
				return "", errors.Wrapf(err, "error running [%s]", ex)
			}
		}

		return controller.FlashAndRedir(true, "Issue solved", returnURL, ps)
	})
}
