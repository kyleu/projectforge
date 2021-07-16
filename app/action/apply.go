package action

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Apply(t Type, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	switch t {
	case TypeCreate:
		return onCreate(cfg, mSvc, pSvc, logger)
	case TypeTest:
		return onTest(cfg, mSvc, pSvc, logger)
	}

	prj, mods, err := prjAndMods(cfg, mSvc, pSvc)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	switch t {
	case TypeBuild:
		return onBuild(prj, cfg, logger)
	case TypeCreate:
		return onCreate(cfg, mSvc, pSvc, logger)
	case TypeDebug:
		return onDebug(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeMerge:
		return onMerge(prj, mods, cfg, mSvc, pSvc, logger)
	case TypePreview:
		return onPreview(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeSlam:
		return onSlam(prj, mods, cfg, mSvc, pSvc, logger)
	case TypeSVG:
		return onSVG(prj, cfg, logger)
	case TypeTest:
		return onTest(cfg, mSvc, pSvc, logger)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t), cfg, logger)
	}
}

func prjAndMods(cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service) (*project.Project, module.Modules, error) {
	path, _ := cfg.GetString("path", true)
	if path == "" {
		path = "."
	}
	prj, err := pSvc.Load(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to load newly created project")
	}
	mods, err := mSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, nil, err
	}
	return prj, mods, nil
}
