package action

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onPreview(prj *project.Project, mods module.Modules, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	start := util.TimerStart()
	_, diffs, err := diffs(prj, mods, true, mSvc, pSvc, logger)
	if err != nil {
		return ret.WithError(err)
	}

	mr := &module.Result{Keys: mods.Keys(), Status: "OK", Diffs: diffs, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
