package action

import (
	"context"

	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func Apply(ctx context.Context, p *Params) *Result {
	if p.Span != nil {
		p.Span.SetAttributes(attribute.String("project", p.ProjectKey), attribute.String("action", p.T.String()))
		defer p.Span.End()
	}

	switch p.T {
	case TypeCreate:
		return onCreate(ctx, p.ProjectKey, p.Cfg, p.RootFiles, p.MSvc, p.PSvc, p.Logger)
	case TypeTest:
		return onTest(ctx, p.Cfg, p.RootFiles, p.MSvc, p.PSvc, p.Logger)
	case TypeDoctor:
		return onDoctor(ctx, p.Cfg, p.PSvc, p.Logger)
	}

	pm, err := getPrjAndMods(ctx, p)
	if err != nil {
		return errorResult(err, p.Cfg, p.Logger)
	}

	switch p.T {
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
		return errorResult(errors.Errorf("invalid action type [%s]", p.T.String()), p.Cfg, p.Logger)
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

func getPrjAndMods(ctx context.Context, p *Params) (*PrjAndMods, error) {
	prj, err := p.PSvc.Get(p.ProjectKey)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load project [%s]", p.ProjectKey)
	}
	mods, err := p.MSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, err
	}
	return &PrjAndMods{Ctx: ctx, Cfg: p.Cfg, Prj: prj, Mods: mods, MSvc: p.MSvc, PSvc: p.PSvc, Logger: p.Logger}, nil
}
