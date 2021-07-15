package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views"
	"github.com/kyleu/projectforge/views/vaction"
	"github.com/valyala/fasthttp"
)

var runContent = util.ValueMap{
	"_":    util.AppName,
	"urls": map[string]string{},
}

func Run(ctx *fasthttp.RequestCtx) {
	act("run", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = runContent
		return render(ctx, as, &views.Debug{}, ps, "actions")
	})
}

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
		switch tgt {
		case "admini":
			cfg["path"] = "../admini"
		case "self":
			cfg["path"] = "."
		case "test":
			cfg.Add("method", actS, "path", "./testproject", "wipe", true)
			actT = action.TypeTest
		default:
			return ersp("invalid target [%s]", tgt)
		}
		result := action.Apply(actT, cfg, controllerModuleSvc, controllerProjectSvc, controllerLogger)
		ps.Data = result
		page := &vaction.Result{Cfg: cfg, Result: result}
		return render(ctx, as, page, ps, "actions", tgt+"-"+actS)
	})
}
