package action

import (
	"context"

	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onDoctor(ctx context.Context, cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	prgs := doctor.Check()
	for _, prg := range prgs {
		ret.AddLog(prg.String())
	}
	return ret
}
