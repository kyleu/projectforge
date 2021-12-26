package action

import (
	"context"

	"github.com/kyleu/projectforge/app/export"
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

	start := util.TimerStart()
	ret := applyBasic(ctx, p)
	if ret == nil {
		if len(p.PSvc.Projects()) == 0 {
			_, err := p.PSvc.Refresh()
			if err != nil {
				return errorResult(err, p.Cfg, p.Logger)
			}
		}
		if p.ProjectKey == "" {
			prj := p.PSvc.ByPath(".")
			p.ProjectKey = prj.Key
		}

		pm, err := getPrjAndMods(ctx, p)
		if err != nil {
			return errorResult(err, p.Cfg, p.Logger)
		}

		ret = applyPrj(ctx, pm, p.T)
	}
	ret.Duration = util.TimerEnd(start)
	return ret
}

func applyBasic(ctx context.Context, p *Params) *Result {
	switch p.T {
	case TypeCreate:
		return onCreate(ctx, p)
	case TypeTest:
		return onTest(ctx, p)
	case TypeDoctor:
		return onDoctor(ctx, p.Cfg, p.PSvc, p.Logger)
	}
	return nil
}

func applyPrj(ctx context.Context, pm *PrjAndMods, t Type) *Result {
	switch t {
	case TypeBuild:
		return onBuild(pm)
	case TypeMerge:
		return onMerge(pm)
	case TypePreview:
		return onPreview(pm)
	case TypeSlam:
		return onSlam(pm)
	case TypeSVG:
		return onSVG(pm)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), pm.Cfg, pm.Logger)
	}
}

type PrjAndMods struct {
	Ctx    context.Context
	Cfg    util.ValueMap
	Prj    *project.Project
	Mods   module.Modules
	MSvc   *module.Service
	PSvc   *project.Service
	ESvc   *export.Service
	Logger *zap.SugaredLogger
}

func getPrjAndMods(ctx context.Context, p *Params) (*PrjAndMods, error) {
	if p.ProjectKey == "" {
		prj := p.PSvc.ByPath("")
		if prj != nil {
			p.ProjectKey = prj.Key
		}
	}

	prj, err := p.PSvc.Get(p.ProjectKey)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load project [%s]", p.ProjectKey)
	}
	mods, err := p.MSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, err
	}
	return &PrjAndMods{Ctx: ctx, Cfg: p.Cfg, Prj: prj, Mods: mods, MSvc: p.MSvc, PSvc: p.PSvc, ESvc: p.ESvc, Logger: p.Logger}, nil
}
