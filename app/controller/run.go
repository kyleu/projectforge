package controller

import (
	"sort"
	"strings"
	"sync"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/telemetry"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vaction"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RunAction(rc *fasthttp.RequestCtx) {
	act("run.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		tgt, err := rcRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		actS, err := rcRequiredString(rc, "act", false)
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
		nc, span := telemetry.StartSpan(ps.Context, "action", "action."+actT.String())
		result := action.Apply(nc, actionParams(span, tgt, actT, cfg, as, ps.Logger))
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
		actS, err := rcRequiredString(rc, "act", false)
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
				nc, span := telemetry.StartSpan(ps.Context, "action", "action."+actT.String())
				result := action.Apply(nc, actionParams(span, p.Key, actT, c, as, ps.Logger))
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
			return strings.ToLower(results[i].Prj.Key) < strings.ToLower(results[j].Prj.Key)
		})
		ps.Data = results
		page := &vaction.Results{T: actT, Ctxs: results}
		return render(rc, as, page, ps, "projects", actT.Title)
	})
}

func actionParams(span trace.Span, tgt string, t action.Type, cfg util.ValueMap, as *app.State, logger *zap.SugaredLogger) *action.Params {
	return &action.Params{
		Span: span, ProjectKey: tgt, T: t, Cfg: cfg,
		MSvc: as.Services.Modules, PSvc: as.Services.Projects, CSvc: as.Services.Codegen, Logger: logger,
	}
}
