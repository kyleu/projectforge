package action

import (
	"context"

	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func onDoctor(ctx context.Context, cfg util.ValueMap, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(TypeDoctor, cfg, logger)
	prjs := pSvc.Projects()
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
