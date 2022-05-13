package action

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Apply(ctx context.Context, p *Params) *Result {
	ctx, span, logger := telemetry.StartSpan(ctx, "action:"+p.T.Key, p.Logger)
	defer span.Complete()
	span.Attribute("project", p.ProjectKey)
	span.Attribute("action", p.T.String())

	timer := util.TimerStart()
	var ret *Result

	defer func() {
		if ret != nil {
			ret.Duration = timer.End()
		}
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				ret = ret.WithError(err)
			} else {
				ret = ret.WithError(errors.New(fmt.Sprint(rec)))
			}
		}
	}()

	ret = applyBasic(ctx, p)
	if ret == nil {
		if len(p.PSvc.Projects()) == 0 {
			_, err := p.PSvc.Refresh()
			if err != nil {
				return errorResult(err, p.T, p.Cfg, logger)
			}
		}
		if p.ProjectKey == "" {
			prj := p.PSvc.ByPath(".")
			p.ProjectKey = prj.Key
		}

		var pm *PrjAndMods
		var err error
		ctx, pm, err = getPrjAndMods(ctx, p)
		if err != nil {
			return errorResult(err, p.T, p.Cfg, logger)
		}

		ret = applyPrj(ctx, pm, p.T)
	}
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
	case TypeAudit:
		return onAudit(ctx, pm)
	case TypeBuild:
		return onBuild(ctx, pm)
	case TypeDebug:
		return onDebug(ctx, pm)
	case TypeMerge:
		return onMerge(ctx, pm)
	case TypePreview:
		return onPreview(ctx, pm)
	case TypeSlam:
		return onSlam(ctx, pm)
	case TypeSVG:
		return onSVG(ctx, pm)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), t, pm.Cfg, pm.Logger)
	}
}

type PrjAndMods struct {
	Cfg    util.ValueMap
	Prj    *project.Project
	Mods   module.Modules
	MSvc   *module.Service
	PSvc   *project.Service
	ESvc   *export.Service
	EArgs  *model.Args
	Logger util.Logger
}

func getPrjAndMods(ctx context.Context, p *Params) (context.Context, *PrjAndMods, error) {
	if p.ProjectKey == "" {
		prj := p.PSvc.ByPath("")
		if prj != nil {
			p.ProjectKey = prj.Key
		}
	}

	prj, err := p.PSvc.Get(p.ProjectKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to load project [%s]", p.ProjectKey)
	}
	if prj.Info != nil {
		_, e := p.MSvc.Register(ctx, prj.Path, prj.Info.ModuleDefs...)
		if e != nil {
			return nil, nil, errors.Wrap(e, "unable to register modules")
		}
	}

	mods, err := p.MSvc.GetModules(prj.Modules...)
	if err != nil {
		return nil, nil, err
	}

	args := &model.Args{}
	if argsX := prj.Info.ModuleArg("export"); argsX != nil {
		err := util.CycleJSON(argsX, args)
		if err != nil {
			return nil, nil, errors.Wrap(err, "export module arguments are invalid")
		}
	}
	args.Modules = mods.Keys()

	pm := &PrjAndMods{Cfg: p.Cfg, Prj: prj, Mods: mods, MSvc: p.MSvc, PSvc: p.PSvc, ESvc: p.ESvc, EArgs: args, Logger: p.Logger}
	return ctx, pm, nil
}
