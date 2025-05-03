package site

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vsite"
)

func featuresMenu(mSvc *module.Service) menu.Items {
	return lo.Map(mSvc.ModulesVisible(), func(m *module.Module, _ int) *menu.Item {
		return &menu.Item{Key: m.Key, Title: m.Title(), Description: m.Description, Icon: m.IconSafe(), Route: m.FeaturesPath()}
	})
}

func featureList(as *app.State, ps *cutil.PageState) (layout.Page, error) {
	mods := as.Services.Modules.ModulesVisible()
	ps.SetTitleAndData("Available Modules", mods)
	return &vsite.FeatureList{Modules: mods}, nil
}

func featureDetail(key string, as *app.State, ps *cutil.PageState) (layout.Page, error) {
	mod, err := as.Services.Modules.Get(key)
	if err != nil {
		return nil, err
	}
	_, html, err := doc.HTMLString("feature:"+mod.Key, []byte(mod.UsageMD), func(s string) (string, string, error) {
		ret, e := cutil.FormatMarkdown(s)
		if e != nil {
			return "", "", e
		}
		if h1Idx := strings.Index(ret, "<h1>"); h1Idx > -1 {
			if h1EndIdx := strings.Index(ret, "</h1>"); h1EndIdx > -1 {
				ret = ret[:h1Idx] + ret[h1EndIdx+5:]
			}
		}
		return "", ret, nil
	})
	if err != nil {
		return nil, err
	}
	ps.SetTitleAndData(mod.Title(), mod)
	return &vsite.FeatureDetail{Module: mod, HTML: html}, nil
}

func featureFiles(path []string, as *app.State, ps *cutil.PageState) ([]string, layout.Page, error) {
	if len(path) < 3 {
		return path, nil, errors.New("invalid path")
	}
	if path[2] != "files" {
		return path, nil, errors.New("invalid file path")
	}
	mod, err := as.Services.Modules.Get(path[1])
	if err != nil {
		return path, nil, err
	}
	u := mod.FeaturesFilePath()
	bc := util.ArrayCopy(path[:2])
	bc = append(bc, "Files||"+u)
	lo.ForEach(path[3:], func(x string, _ int) {
		u += "/" + x
		bc = append(bc, x+"||"+u)
	})
	ps.SetTitleAndData("Module ["+mod.Title()+"]", mod)
	ps.Data = mod
	return bc, &vsite.FeatureFiles{Module: mod, Path: path[3:]}, nil
}
