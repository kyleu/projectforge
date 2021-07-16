package action

import (
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func onMerge(prj *project.Project, mods module.Modules, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	var res module.Results
	for _, mod := range mods {
		r, err := merge(prj, mod, mSvc, pSvc)
		if err != nil {
			return ret.WithError(err)
		}
		res = append(res, r)
		logger.Info("applied module [" + mod.Key + "]")
	}
	ret.Modules = res
	return ret
}

func merge(prj *project.Project, mod *module.Module, mSvc *module.Service, pSvc *project.Service) (*module.Result, error) {
	start := util.TimerStart()
	_, diffs, err := diffs(prj, mod, mSvc, pSvc, true)
	if err != nil {
		return nil, err
	}

	for _, f := range diffs {
		switch f.Status {
		case file.StatusIdentical, file.StatusMissing, file.StatusSkipped:
			// noop
		case file.StatusDifferent, file.StatusNew:
			// noop
		default:
			return nil, errors.Errorf("unhandled diff status [%s]", f.Status)
		}
	}

	return &module.Result{Key: mod.Key, Status: "OK", Diffs: diffs, Duration: util.TimerEnd(start)}, nil
}
