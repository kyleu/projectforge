package action

import (
	"context"

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
	for _, r := range res {
		ret.AddLog("%s: %s", r.Title, r.Status)
		for _, l := range r.Logs {
			ret.AddLog(" - %s", l)
		}
		for _, e := range r.Errors {
			ret.AddWarn(" - %s", e.String())
		}
		for _, s := range r.Solution {
			ret.AddDebug(" - %s", s)
		}
	}
	return ret
}
