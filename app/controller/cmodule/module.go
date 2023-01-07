package cmodule

import (
	"strings"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/vmodule"
)

func ModuleList(rc *fasthttp.RequestCtx) {
	controller.Act("module.root", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mods := as.Services.Modules.Modules()
		ps.Title = "Module Listing"
		ps.Data = mods
		return controller.Render(rc, as, &vmodule.List{Modules: mods}, ps, "modules")
	})
}

func ModuleDetail(rc *fasthttp.RequestCtx) {
	controller.Act("module.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mod, err := getModule(rc, as, ps)
		if err != nil {
			return "", err
		}
		var usages project.Projects
		for _, p := range as.Services.Projects.Projects() {
			if p.HasModule(mod.Key) {
				usages = append(usages, p)
			}
		}
		_, html, err := doc.HTMLString("module:"+mod.Key, []byte(mod.UsageMD), func(s string) (string, string, error) {
			ret, err := cutil.FormatMarkdown(s)
			if err != nil {
				return "", "", err
			}
			if aIdx := strings.Index(ret, "CC0</a>"); aIdx > -1 {
				ret = strings.TrimPrefix(ret[aIdx+7:], "</p>")
			}
			return "", ret, nil
		})
		if err != nil {
			return "", err
		}
		ps.Data = mod
		ps.Title = mod.Title()
		return controller.Render(rc, as, &vmodule.Detail{Module: mod, HTML: html, Usages: usages}, ps, "modules", mod.Key)
	})
}

func getModule(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*module.Module, error) {
	key, err := cutil.RCRequiredString(rc, "key", true)
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
