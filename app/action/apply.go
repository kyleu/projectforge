package action

import (
	"context"

	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Apply(ctx context.Context, span trace.Span, projectKey string, t Type, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	if span != nil {
		span.SetAttributes(attribute.String("project", projectKey), attribute.String("action", t.String()))
		defer span.End()
	}

	switch t {
	case TypeCreate:
		return onCreate(ctx, projectKey, cfg, mSvc, pSvc, logger)
	case TypeTest:
		return onTest(ctx, cfg, mSvc, pSvc, logger)
	case TypeDoctor:
		return onDoctor(ctx, cfg, pSvc, logger)
	}

	pm, err := getPrjAndMods(ctx, projectKey, cfg, mSvc, pSvc, logger)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	switch t {
	case TypeBuild:
		return onBuild(pm)
	case TypeDebug:
		return onDebug(pm)
	case TypeMerge:
		return onMerge(pm)
	case TypePreview:
		return onPreview(pm)
	case TypeSlam:
		return onSlam(pm)
	case TypeSVG:
		return onSVG(pm)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), cfg, logger)
	}
}

type PrjAndMods struct {
	Ctx    context.Context
	Cfg    util.ValueMap
	Prj    *project.Project
	Mods   module.Modules
	MSvc   *module.Service
	PSvc   *project.Service
	Logger *zap.SugaredLogger
}

func getPrjAndMods(ctx context.Context, key string, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) (*PrjAndMods, error) {
	prj, err := pSvc.Get(key)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load project [%s]", key)
	}
	mods, err := mSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, err
	}
	return &PrjAndMods{Ctx: ctx, Cfg: cfg, Prj: prj, Mods: mods, MSvc: mSvc, PSvc: pSvc, Logger: logger}, nil
}
