package controller

import (
	"sort"
	"sync"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/telemetry"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vaction"
	"github.com/valyala/fasthttp"
)

func RunAction(ctx *fasthttp.RequestCtx) {
	act("run.action", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		tgt, err := ctxRequiredString(ctx, "tgt", false)
		if err != nil {
			return "", err
		}
		actS, err := ctxRequiredString(ctx, "act", false)
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
		ctx.QueryArgs().VisitAll(func(k []byte, v []byte) {
			cfg[string(k)] = string(v)
		})
		nc, span := telemetry.StartSpan(ps.Context, "action", "action." + actT.String())
		result := action.Apply(nc, span, tgt, actT, cfg, as.Services.Modules, as.Services.Projects, ps.Logger)
		ps.Data = result
		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}}
		return render(ctx, as, page, ps, "projects", prj.Key, actT.Title)
	})
}

func RunAllActions(ctx *fasthttp.RequestCtx) {
	act("run.all.actions", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := ctxRequiredString(ctx, "act", false)
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
				result := action.Apply(nc, span, p.Key, actT, c, as.Services.Modules, as.Services.Projects, ps.Logger)
				rc := &action.ResultContext{Prj: p, Cfg: c, Res: result}
				mutex.Lock()
				results = append(results, rc)
				wg.Done()
				mutex.Unlock()
			}()
		}
		wg.Wait()
		sort.Slice(results, func(i, j int) bool {
			return results[i].Prj.Key < results[j].Prj.Key
		})
		ps.Data = results
		page := &vaction.Results{T: actT, Ctxs: results}
		return render(ctx, as, page, ps, "projects", actT.Title)
	})
}
