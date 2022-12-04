package site

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vsite"
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
	html, err := doc.HTML("components/" + key + ".md")
	if err != nil {
		return "", "", err
	}
	title := util.StringToTitle(key)
	if h1Idx := strings.Index(html, "<h1>"); h1Idx > -1 {
		if h1EndIdx := strings.Index(html, "</h1>"); h1EndIdx > -1 {
			title = html[h1Idx+4 : h1EndIdx]
		}
		println(icon)
		ic := fmt.Sprintf(`<svg class="icon" style="width: 36px; height: 36px;"><use xlink:href="#svg-%s" /></svg> `, icon)
		html = html[:h1Idx+4] + ic + html[h1Idx+4:]
	}
	return title, html, nil
}
