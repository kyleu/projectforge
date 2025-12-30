package cmodule

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
		fsys := as.Services.Modules.GetFilesystem(mod.Key)
		ps.SetTitleAndData(fmt.Sprintf("[%s] Files", mod.Key), mod)
		return controller.Render(r, as, &vmodule.Files{Module: mod, FS: fsys}, ps, "modules", mod.Key, "Files**folder")
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
		fsys := as.Services.Modules.GetFilesystem(mod.Key)
		path := util.StringSplitAndTrim(pathS, "/")
		if cutil.QueryStringString(ps.URI, "download") == "true" {
			b, err := fsys.ReadFile(util.StringFilePath(path...))
			if err != nil {
				return "", errors.Wrapf(err, "unable to read module file [%s] for download", pathS)
			}
			return cutil.RespondDownload(path[len(path)-1], b, ps.W)
		}
		bcAppend := "||/m/" + mod.Key + "/fs"
		bc := []string{"modules", mod.Key, "Files" + bcAppend + "**folder"}
		lo.ForEach(path, func(x string, _ int) {
			bcAppend += "/" + x
			b := x + bcAppend
			bc = append(bc, b)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] /%s", mod.Key, util.StringJoin(path, "/")), pathS)
		return controller.Render(r, as, &vmodule.Files{Module: mod, Path: path, FS: fsys}, ps, bc...)
	})
}
