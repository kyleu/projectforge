package cmodule

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleFileRoot(rc *fasthttp.RequestCtx) {
	controller.Act("module.file.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Files", mod.Key), mod)
		return controller.Render(rc, as, &vmodule.Files{Module: mod}, ps, "modules", mod.Key, "Files**folder")
	})
}

func ModuleFile(rc *fasthttp.RequestCtx) {
	controller.Act("module.file", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
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
		bc := []string{"modules", mod.Key, "Files" + bcAppend + "**folder"}
		lo.ForEach(path, func(x string, _ int) {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] /%s", mod.Key, strings.Join(path, "/")), pathS)
		return controller.Render(rc, as, &vmodule.Files{Module: mod, Path: path}, ps, bc...)
	})
}
