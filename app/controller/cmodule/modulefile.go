package cmodule

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleFileRoot(w http.ResponseWriter, r *http.Request) {
	controller.Act("module.file.root", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Files", mod.Key), mod)
		return controller.Render(w, r, as, &vmodule.Files{Module: mod}, ps, "modules", mod.Key, "Files**folder")
	})
}

func ModuleFile(w http.ResponseWriter, r *http.Request) {
	controller.Act("module.file", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(r, as, ps)
		if err != nil {
			return "", err
		}

		pathS, err := cutil.PathString(r, "path", false)
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
		return controller.Render(w, r, as, &vmodule.Files{Module: mod, Path: path}, ps, bc...)
	})
}
