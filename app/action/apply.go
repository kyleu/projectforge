package action

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Apply(projectKey string, t Type, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	switch t {
	case TypeCreate:
		return onCreate(projectKey, cfg, mSvc, pSvc, logger)
	case TypeTest:
		return onTest(cfg, mSvc, pSvc, logger)
	case TypeDoctor:
		return onDoctor(cfg, logger)
	}

	prj, mods, err := prjAndMods(projectKey, mSvc, pSvc)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	switch t {
	case TypeBuild:
		return onBuild(prj, cfg, logger)
	case TypeDebug:
		return onDebug(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeMerge:
		return onMerge(prj, mods, cfg, mSvc, pSvc, logger)
	case TypePreview:
		return onPreview(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeSlam:
		return onSlam(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeSVG:
		return onSVG(prj, cfg, pSvc, logger)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), cfg, logger)
	}
}

func prjAndMods(key string, mSvc *module.Service, pSvc *project.Service) (*project.Project, module.Modules, error) {
	prj, err := pSvc.Get(key)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to load project [%s]", key)
	}
	mods, err := mSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, nil, err
	}
	return prj, mods, nil
}
