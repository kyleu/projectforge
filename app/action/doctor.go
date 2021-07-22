package action

import (
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onDoctor(cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	ret.AddLog("Doctor doctor, gimme the news")
	return ret
}
