package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
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
		mod, err := getModule(ctx, as, ps)
		if err != nil {
			return "", err
		}
		return render(ctx, as, &vmodule.Detail{Module: mod}, ps, "modules", mod.Key)
	})
}

func ModuleFileRoot(ctx *fasthttp.RequestCtx) {
	act("module.file.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(ctx, as, ps)
		if err != nil {
			return "", err
		}
		return render(ctx, as, &vmodule.Files{Module: mod}, ps, "modules", mod.Key, "Files")
	})
}

func ModuleFile(ctx *fasthttp.RequestCtx) {
	act("module.file", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(ctx, as, ps)
		if err != nil {
			return "", err
		}

		pathS, err := ctxRequiredString(ctx, "path", false)
		if err != nil {
			return "", err
		}
		path := util.SplitAndTrim(pathS, "/")
		bcAppend := "||/m/" + mod.Key + "/fs"
		bc := []string{"modules", mod.Key, "Files" + bcAppend}
		for _, x := range path {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		}
		return render(ctx, as, &vmodule.Files{Module: mod, Path: path}, ps, bc...)
	})
}

func getModule(ctx *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*module.Module, error) {
	key, err := ctxRequiredString(ctx, "key", true)
	if err != nil {
		return nil, err
	}

	mod, err := as.Services.Modules.Get(key)
	if err != nil {
		return nil, err
	}

	ps.Title = fmt.Sprintf("%s", mod.Title())
	ps.Data = mod

	return mod, nil
}
