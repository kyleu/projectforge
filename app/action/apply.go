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

	start := util.TimerStart()

	var ret *Result
	switch p.T {
	case TypeCreate:
		ret = onCreate(ctx, p)
	case TypeTest:
		ret = onTest(ctx, p)
	case TypeDoctor:
		ret = onDoctor(ctx, p.Cfg, p.PSvc, p.Logger)
	}

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

	switch p.T {
	case TypeBuild:
		ret = onBuild(pm, p.Cfg)
	case TypeMerge:
		ret = onMerge(pm)
	case TypePreview:
		ret = onPreview(pm)
	case TypeSlam:
		ret = onSlam(pm)
	case TypeSVG:
		ret = onSVG(pm)
	default:
		ret = errorResult(errors.Errorf("invalid action type [%s]", p.T.String()), p.Cfg, p.Logger)
	}

	ret.Duration = util.TimerEnd(start)
	return ret
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
	return &PrjAndMods{Ctx: ctx, Cfg: p.Cfg, Prj: prj, Mods: mods, MSvc: p.MSvc, PSvc: p.PSvc, Logger: p.Logger}, nil
}
