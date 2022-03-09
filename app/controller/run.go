package controller

import (
	"sort"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vaction"
)

func RunAction(rc *fasthttp.RequestCtx) {
	act("run.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		tgt, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		actS, err := RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}

		cfg := util.ValueMap{}
		actT := action.TypeFromString(actS)
		prj, err := as.Services.Projects.Get(tgt)
		if err != nil {
			return "", err
		}
		cfg["path"] = prj.Path
		rc.QueryArgs().VisitAll(func(k []byte, v []byte) {
			cfg[string(k)] = string(v)
		})
		nc, span, logger := telemetry.StartSpan(ps.Context, "action:"+actT.String(), ps.Logger)
		result := action.Apply(nc, actionParams(span, tgt, actT, cfg, as, logger))
		if result.Project == nil {
			result.Project = prj
		}
		ps.Data = result
		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}}
		return render(rc, as, page, ps, "projects", prj.Key, actT.Title)
	})
}

func RunAllActions(rc *fasthttp.RequestCtx) {
	act("run.all.actions", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		cfg := util.ValueMap{}
		actT := action.TypeFromString(actS)
		prjs := as.Services.Projects.Projects()

		var results action.ResultContexts
		var mutex sync.Mutex
		var wg sync.WaitGroup
		wg.Add(len(prjs))
		for _, prj := range prjs {
			p := prj
			go func() {
				c := cfg.Clone()
				c["path"] = p.Path
				nc, span, logger := telemetry.StartSpan(ps.Context, "action:"+actT.String(), ps.Logger)
				result := action.Apply(nc, actionParams(span, p.Key, actT, c, as, logger))
				if result.Project == nil {
					result.Project = p
				}
				rc := &action.ResultContext{Prj: p, Cfg: c, Res: result}
				mutex.Lock()
				results = append(results, rc)
				wg.Done()
				mutex.Unlock()
			}()
		}
		wg.Wait()
		sort.Slice(results, func(i int, j int) bool {
			return strings.ToLower(results[i].Prj.Title()) < strings.ToLower(results[j].Prj.Title())
		})
		ps.Data = results
		page := &vaction.Results{T: actT, Ctxs: results}
		return render(rc, as, page, ps, "projects", actT.Title)
	})
}

func actionParams(span *telemetry.Span, tgt string, t action.Type, cfg util.ValueMap, as *app.State, logger *zap.SugaredLogger) *action.Params {
	return &action.Params{
		Span: span, ProjectKey: tgt, T: t, Cfg: cfg,
		MSvc: as.Services.Modules, PSvc: as.Services.Projects, ESvc: as.Services.Export, Logger: logger,
	}
}
