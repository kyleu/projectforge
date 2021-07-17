package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/views/vmodule"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ModuleList(ctx *fasthttp.RequestCtx) {
	act("module.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		mods := as.Services.Modules.Modules()
		ps.Title = "Module Listing"
		ps.Data = mods
		return render(ctx, as, &vmodule.List{Modules: mods}, ps, "modules")
	})
}

func ModuleDetail(ctx *fasthttp.RequestCtx) {
	act("module.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", true)
		if err != nil {
			return "", err
		}

		mod, err := as.Services.Modules.Get(key)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (module %s)", mod.Title(), mod.Key)
		ps.Data = mod
		return render(ctx, as, &vmodule.Detail{Module: mod}, ps, "modules", mod.Key)
	})
}
