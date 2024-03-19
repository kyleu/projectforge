package cmodule

import (
	"net/http"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleList(w http.ResponseWriter, r *http.Request) {
	controller.Act("module.root", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mods := as.Services.Modules.Modules()
		ps.SetTitleAndData("Module Listing", mods)
		dir := as.Services.Modules.ConfigDirectory()
		return controller.Render(w, r, as, &vmodule.List{Modules: mods, Dir: dir}, ps, "modules")
	})
}

func ModuleDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("module.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(r, as, ps)
		if err != nil {
			return "", err
		}
		usages := lo.Filter(as.Services.Projects.Projects(), func(p *project.Project, _ int) bool {
			return p.HasModule(mod.Key)
		})
		_, html, err := doc.HTMLString("module:"+mod.Key, []byte(mod.UsageMD), func(s string) (string, string, error) {
			ret, e := cutil.FormatMarkdown(s)
			if e != nil {
				return "", "", e
			}
			if aIdx := strings.Index(ret, "CC0</a>"); aIdx > -1 {
				ret = strings.TrimPrefix(ret[aIdx+7:], "</p>")
			}
			return "", ret, nil
		})
		if err != nil {
			return "", err
		}
		dir := mod.Files.Root()
		ps.SetTitleAndData(mod.Title(), mod)
		return controller.Render(w, r, as, &vmodule.Detail{Module: mod, HTML: html, Usages: usages, Dir: dir}, ps, "modules", mod.Key)
	})
}

func getModule(r *http.Request, as *app.State, ps *cutil.PageState) (*module.Module, error) {
	key, err := cutil.RCRequiredString(r, "key", true)
	if err != nil {
		return nil, err
	}
	mod, err := as.Services.Modules.Get(key)
	if err != nil {
		return nil, err
	}
	ps.SetTitleAndData(mod.Title(), mod)
	return mod, nil
}
