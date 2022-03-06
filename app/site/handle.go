// Package site $PF_IGNORE$
package site

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"projectforge.dev/views/verror"

	"projectforge.dev/app"
	"projectforge.dev/app/controller/cutil"
	"projectforge.dev/app/site/download"
	"projectforge.dev/app/util"
	"projectforge.dev/doc"
	"projectforge.dev/views/layout"
	"projectforge.dev/views/vsite"
)

func Handle(path []string, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, layout.Page, []string, error) {
	if len(path) == 0 {
		msg := "\n  <meta name=\"go-source\" content=\"projectforge.dev %s %s/tree/master{/dir} %s/blob/master{/dir}/{file}#L{line}\">"
		ps.HeaderContent = fmt.Sprintf(msg, util.AppSource, util.AppSource, util.AppSource)
		ps.Data = siteData("Welcome to the marketing site!")
		return "", &vsite.Index{}, path, nil
	}

	var page layout.Page
	var err error
	bc := path
	switch path[0] {
	case keyFeatures:
		switch {
		case len(path) == 1:
			page, err = featureList(as, ps)
		case len(path) == 2:
			page, err = featureDetail(path[1], as, ps)
		default:
			bc, page, err = featureFiles(path, as, ps)
		}
	case keyDownload:
		dls := download.GetLinks(as.BuildInfo.Version)
		ps.Data = map[string]interface{}{"base": "https://github.com/kyleu/projectforge/releases/download/v" + as.BuildInfo.Version, "links": dls}
		page = &vsite.Download{Links: dls}
	case keyInstall:
		page, err = mdTemplate("Installation", "This static page contains installation instructions", "installation.md", "code", ps)
	case keyContrib:
		page, err = mdTemplate("Contributing", "This static page describes how to build "+util.AppName, "contributing.md", "cog", ps)
	case keyTech:
		page, err = mdTemplate("Technology", "This static page describes the technology used in "+util.AppName, "technology.md", "shield", ps)
	default:
		page, err = mdTemplate("Documentation", "Documentation for "+util.AppName, path[0]+".md", "", ps)
		if err != nil {
			page = &verror.NotFound{Path: "/" + strings.Join(path, "/")}
			err = nil
		}
	}
	return "", page, bc, err
}

func siteData(result string, kvs ...string) map[string]interface{} {
	ret := map[string]interface{}{"app": util.AppName, "url": util.AppURL, "result": result}
	for i := 0; i < len(kvs); i += 2 {
		ret[kvs[i]] = kvs[i+1]
	}
	return ret
}

func mdTemplate(title string, description string, path string, icon string, ps *cutil.PageState) (layout.Page, error) {
	if icon == "" {
		icon = "cog"
	}
	ps.Data = siteData(title, "description", description)
	ps.Title = title
	html, err := doc.HTML(path)
	if err != nil {
		return nil, err
	}
	if h1Idx := strings.Index(html, "<h1>"); h1Idx > -1 {
		ic := fmt.Sprintf(`<svg class="icon" style="width: 36px; height: 36px;"><use xlink:href="#svg-%s" /></svg> `, icon)
		html = html[:h1Idx+4] + ic + html[h1Idx+4:]
	}
	page := &vsite.MarkdownPage{Title: title, HTML: html}
	return page, nil
}
