package action

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ApplyAll(ctx context.Context, prjs project.Projects, actT Type, cfg util.ValueMap, as *app.State, logger util.Logger) ResultContexts {
	serial := cfg.GetBoolOpt("serial") || cfg.GetStringOpt("mode") == refreshMode
	mu := sync.Mutex{}
	nonSelf := prjs.WithoutTags("self")
	results, _ := util.AsyncCollect(nonSelf, func(prj *project.Project) (*ResultContext, error) {
		if serial {
			mu.Lock()
			defer mu.Unlock()
		}
		return funcName(prj, cfg, actT, as, logger, ctx)
	})
	if len(prjs) > len(nonSelf) {
		for _, self := range prjs.WithTags("self") {
			x, _ := funcName(self, cfg, actT, as, logger, ctx)
			results = append(results, x)
		}
	}
	slices.SortFunc(results, func(l *ResultContext, r *ResultContext) int {
		return cmp.Compare(strings.ToLower(l.Prj.Title()), strings.ToLower(r.Prj.Title()))
	})
	return results
}

func funcName(prj *project.Project, cfg util.ValueMap, actT Type, as *app.State, logger util.Logger, ctx context.Context) (*ResultContext, error) {
	mSvc, pSvc, eSvc, xSvc := as.Services.Modules, as.Services.Projects, as.Services.Export, as.Services.Exec
	c := cfg.Clone()
	prms := &Params{ProjectKey: prj.Key, T: actT, Cfg: cfg, MSvc: mSvc, PSvc: pSvc, XSvc: xSvc, ESvc: eSvc, Logger: logger}
	result := Apply(ctx, prms)
	if result.Project == nil {
		result.Project = prj
	}
	return &ResultContext{Prj: prj, Cfg: c, Res: result}, nil
}

func Apply(ctx context.Context, p *Params) (ret *Result) {
	ctx, span, logger := telemetry.StartSpan(ctx, "action:"+p.T.Key, p.Logger)
	defer span.Complete()
	span.Attribute("project", p.ProjectKey)
	span.Attribute("action", p.T.String())

	timer := util.TimerStart()

	defer func() {
		if ret == nil {
			ret = &Result{}
		}
		ret.Duration = timer.End()
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				ret = ret.WithError(err)
			} else {
				ret = ret.WithError(errors.New(fmt.Sprint(rec)))
			}
		}
	}()

	ret = applyBasic(ctx, p)
	if ret != nil {
		return ret
	}
	if len(p.PSvc.Projects()) == 0 {
		_, err := p.PSvc.Refresh(p.Logger)
		if err != nil {
			return errorResult(err, p.T, p.Cfg, logger)
		}
	}
	prj := p.PSvc.Default()
	if p.ProjectKey == "" {
		p.ProjectKey = prj.Key
	}
	if prj == nil || strings.Contains(prj.Name, "(missing)") {
		return errorResult(errors.New("no project found in current directory"), p.T, p.Cfg, logger)
	}

	var pm *PrjAndMods
	var err error
	ctx, pm, err = getPrjAndMods(ctx, p)
	if err != nil {
		return errorResult(err, p.T, p.Cfg, logger)
	}
	return applyPrj(ctx, pm, p.T)
}

func applyBasic(ctx context.Context, p *Params) *Result {
	switch p.T {
	case TypeCreate:
		return onCreate(ctx, p)
	case TypeTest:
		return onTest(ctx, p)
	case TypeDoctor:
		return onDoctor(ctx, p.Cfg, p.PSvc, p.MSvc, p.Logger)
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
		return onDebug(pm)
	case TypeGenerate:
		return onGenerate(pm)
	case TypePreview:
		return onPreview(pm)
	case TypeRules:
		return onRules(pm)
	case TypeSVG:
		return onSVG(ctx, pm)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), t, pm.Cfg, pm.Logger)
	}
}
