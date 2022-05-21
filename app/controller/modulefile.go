package controller

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleFileRoot(rc *fasthttp.RequestCtx) {
	act("module.file.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.Title = fmt.Sprintf("[%s] Files", mod.Key)
		return render(rc, as, &vmodule.Files{Module: mod}, ps, "modules", mod.Key, "Files")
	})
}

func ModuleFile(rc *fasthttp.RequestCtx) {
	act("module.file", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}

		pathS, err := cutil.RCRequiredString(rc, "path", false)
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
		ps.Title = fmt.Sprintf("[%s] /%s", mod.Key, strings.Join(path, "/"))
		return render(rc, as, &vmodule.Files{Module: mod, Path: path}, ps, bc...)
	})
}
