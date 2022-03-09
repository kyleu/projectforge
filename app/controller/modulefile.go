package controller

import (
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vmodule"

	"projectforge.dev/projectforge/app/controller/cutil"

	"projectforge.dev/projectforge/app"
)

func ModuleFileRoot(rc *fasthttp.RequestCtx) {
	act("module.file.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}
		return render(rc, as, &vmodule.Files{Module: mod}, ps, "modules", mod.Key, "Files")
	})
}

func ModuleFile(rc *fasthttp.RequestCtx) {
	act("module.file", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}

		pathS, err := RCRequiredString(rc, "path", false)
		if err != nil {
			return "", err
		}
		path := util.StringSplitAndTrim(pathS, "/")
		bcAppend := "||/m/" + mod.Key + "/fs"
		bc := []string{"modules", mod.Key, "Files" + bcAppend}
		for _, x := range path {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		}
		return render(rc, as, &vmodule.Files{Module: mod, Path: path}, ps, bc...)
	})
}
