package site

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/menu"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/views/layout"
	"github.com/kyleu/projectforge/views/vsite"
)

func featuresMenu(as *app.State) menu.Items {
	if as.Services == nil {
		as.Services = &app.Services{
			Modules:  module.NewService(as.Files, as.Logger),
		}
	}
	ms := as.Services.Modules.Modules()
	ret := make(menu.Items, 0, len(ms))
	for _, m := range ms {
		ret = append(ret, &menu.Item{Key: m.Key, Title: m.Title(), Description: m.Description, Icon: m.Icon, Route: "/features/" + m.Key})
	}
	return ret
}

func featureList(as *app.State, ps *cutil.PageState) (layout.Page, error) {
	mods := as.Services.Modules.Modules()
	ps.Data = mods
	return &vsite.FeatureList{Modules: mods}, nil
}

func featureDetail(key string, as *app.State, ps *cutil.PageState) (layout.Page, error) {
	mod, err := as.Services.Modules.Get(key)
	if err != nil {
		return nil, err
	}
	ps.Data = mod
	return &vsite.FeatureDetail{Module: mod}, nil
}
