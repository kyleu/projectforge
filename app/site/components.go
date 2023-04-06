package site

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vsite"
)

func componentsMenu(logger util.Logger) menu.Items {
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

func componentList(ps *cutil.PageState) (layout.Page, error) {
	ps.Title = "Available Components"
	ps.Data = vsite.AllComponents
	return &vsite.ComponentList{}, nil
}

func componentDetail(key string, ps *cutil.PageState) (layout.Page, error) {
	c := vsite.AllComponents.Get(key)
	if c == nil {
		return nil, errors.Errorf("no component available with key [%s]", key)
	}
	ps.AddIcon(c.Icon)
	ps.Title = c.Title
	ps.Data = c
	return &vsite.ComponentDetail{Component: c}, nil
}

func componentTemplate(key string) (string, string, error) {
	title, html, err := doc.HTML("components:"+key, "components/"+key+".md", func(s string) (string, string, error) {
		ret, err := cutil.FormatMarkdown(s)
		if err != nil {
			return "", "", err
		}
		title := util.StringToTitle(key)
		if h1Idx := strings.Index(ret, "<h1>"); h1Idx > -1 {
			if h1EndIdx := strings.Index(ret, "</h1>"); h1EndIdx > -1 {
				title = ret[h1Idx+4 : h1EndIdx]
				ret = ret[:h1Idx] + ret[h1EndIdx+5:]
			}
		}
		return title, ret, nil
	})
	if err != nil {
		return "", "", err
	}
	return title, html, nil
}
