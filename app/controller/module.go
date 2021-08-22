package controller

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/views/vmodule"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ModuleList(rc *fasthttp.RequestCtx) {
	act("module.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mods := as.Services.Modules.Modules()
		ps.Title = "Module Listing"
		ps.Data = mods
		return render(rc, as, &vmodule.List{Modules: mods}, ps, "modules")
	})
}

func ModuleDetail(rc *fasthttp.RequestCtx) {
	act("module.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}
		return render(rc, as, &vmodule.Detail{Module: mod}, ps, "modules", mod.Key)
	})
}

func getModule(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*module.Module, error) {
	key, err := rcRequiredString(rc, "key", true)
	if err != nil {
		return nil, err
	}

	mod, err := as.Services.Modules.Get(key)
	if err != nil {
		return nil, err
	}

	ps.Title = mod.Title()
	ps.Data = mod

	return mod, nil
}
