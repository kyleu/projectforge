package action

import (
	"context"

	"github.com/kyleu/projectforge/app/doctor/checks"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onDoctor(ctx context.Context, cfg util.ValueMap, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	prjs := pSvc.Projects()
	res := checks.CheckAll(prjs.AllModules())
	for _, prg := range res {
		ret.AddLog(prg.String())
	}
	return ret
}
