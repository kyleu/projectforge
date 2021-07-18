package controller

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vmodule"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

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
