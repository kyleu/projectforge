package action

import (
	"context"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func onDoctor(ctx context.Context, cfg util.ValueMap, pSvc *project.Service, mSvc *module.Service, logger util.Logger) *Result {
	ret := newResult(TypeDoctor, nil, cfg, logger)
	prjs := pSvc.Projects()
	checks.CurrentModuleDeps = mSvc.Deps()
	res := checks.CheckAll(ctx, prjs.AllModules(), logger)
	lo.ForEach(res, func(r *doctor.Result, _ int) {
		ret.AddLog("%s: %s", r.Title, r.Status)
		lo.ForEach(r.Logs, func(l string, _ int) {
			ret.AddLog(" - %s", l)
		})
		lo.ForEach(r.Errors, func(e *doctor.Error, _ int) {
			ret.AddWarn(" - %s", e.String())
		})
		lo.ForEach(r.CleanSolutions(), func(s string, _ int) {
			ret.AddDebug(" - %s", s)
		})
	})
	return ret
}
