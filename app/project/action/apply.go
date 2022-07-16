package action

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ApplyAll(ctx context.Context, prjs project.Projects, actT Type, cfg util.ValueMap, as *app.State, logger util.Logger) []*ResultContext {
	serial := cfg.GetBoolOpt("serial") || cfg.GetStringOpt("mode") == refreshMode
	mu := sync.Mutex{}
	mSvc, pSvc, eSvc := as.Services.Modules, as.Services.Projects, as.Services.Export
	results, _ := util.AsyncCollect(prjs, func(prj *project.Project) (*ResultContext, error) {
		if serial {
			mu.Lock()
			defer mu.Unlock()
		}
		c := cfg.Clone()
		prms := &Params{ProjectKey: prj.Key, T: actT, Cfg: cfg, MSvc: mSvc, PSvc: pSvc, ESvc: eSvc, Logger: logger}
		result := Apply(ctx, prms)
		if result.Project == nil {
			result.Project = prj
		}
		return &ResultContext{Prj: prj, Cfg: c, Res: result}, nil
	})
	slices.SortFunc(results, func(l *ResultContext, r *ResultContext) bool {
		return strings.ToLower(l.Prj.Title()) < strings.ToLower(r.Prj.Title())
	})
	return results
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
	if ret == nil {
		if len(p.PSvc.Projects()) == 0 {
			_, err := p.PSvc.Refresh(p.Logger)
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
		return onDebug(ctx, pm)
	case TypeGenerate:
		return onGenerate(ctx, pm)
	case TypePreview:
		return onPreview(ctx, pm)
	case TypeSVG:
		return onSVG(ctx, pm)
	default:
		return errorResult(errors.Errorf("invalid action type [%s]", t.String()), t, pm.Cfg, pm.Logger)
	}
}
