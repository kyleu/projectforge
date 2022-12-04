package site

import (
	"context"
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vsite"
	"strings"
)

func componentsMenu(ctx context.Context, logger util.Logger) menu.Items {
	if vsite.AllComponents == nil {
		var err error
		vsite.AllComponents, err = loadComponents()
		if err != nil {
			logger.Warn(err)
		}
	}
	ret := make(menu.Items, 0, len(vsite.AllComponents))
	for _, c := range vsite.AllComponents {
		ret = append(ret, &menu.Item{Key: c.Key, Title: c.Title, Description: c.Description, Icon: c.Icon, Route: "/components/" + c.Key})
	}
	return ret
}

func componentList(as *app.State, ps *cutil.PageState) (layout.Page, error) {
	ps.Title = "Available Components"
	ps.Data = vsite.AllComponents
	return &vsite.ComponentList{}, nil
}

func componentDetail(key string, as *app.State, ps *cutil.PageState) (layout.Page, error) {
	c := vsite.AllComponents.Get(key)
	if c == nil {
		return nil, errors.Errorf("no component available with key [%s]", key)
	}
	ps.AddIcon(c.Icon)
	return &vsite.ComponentDetail{Component: c}, nil
}

func componentTemplate(key string, icon string) (string, string, error) {
	html, err := doc.HTML("components/"+key+".md", func(s string) (string, error) {
		ret, err := cutil.FormatMarkdown(s)
		if err != nil {
			return "", err
		}
		if h1Idx := strings.Index(ret, "<h1>"); h1Idx > -1 {
			if h1EndIdx := strings.Index(ret, "</h1>"); h1EndIdx > -1 {
				ret = ret[:h1Idx] + ret[h1EndIdx+5:]
			}
		}
		return ret, nil
	})
	if err != nil {
		return "", "", err
	}
	title := util.StringToTitle(key)
	return title, html, nil
}
