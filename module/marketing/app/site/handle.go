package site

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site/download"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/doc"
	"{{{ .Package }}}/views"
	"{{{ .Package }}}/views/layout"
	"{{{ .Package }}}/views/verror"
	"{{{ .Package }}}/views/vsite"
)

func Handle(path []string, rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, layout.Page, []string, error) {
	if len(path) == 0 {
		ps.Data = siteData("Welcome to the marketing site!")
		return "", &vsite.Index{}, path, nil
	}

	var page layout.Page
	var err error
	switch path[0] {
	case util.AppKey:
		msg := "\n  " +
			"<meta name=\"go-import\" content=\"{{{ .Package }}} git %s\">\n  " +
			"<meta name=\"go-source\" content=\"{{{ .Package }}} %s %s/tree/master{/dir} %s/blob/master{/dir}/{file}#L{line}\">"
		ps.HeaderContent = fmt.Sprintf(msg, util.AppSource, util.AppSource, util.AppSource, util.AppSource)
		return "", &vsite.GoSource{}, path, nil
	case keyAbout:
		ps.Title = "About " + util.AppName
		ps.Data = util.AppName + " v" + as.BuildInfo.Version
		page = &views.About{}
	case keyDownload:
		dls := download.GetLinks(as.BuildInfo.Version)
		ps.Title = "Downloads"
		ps.Data = util.ValueMap{"base": "https://{{{ .Package }}}/releases/download/v" + as.BuildInfo.Version, "links": dls}
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
	return "", page, path, err
}

func siteData(result string, kvs ...string) util.ValueMap {
	ret := util.ValueMap{"app": util.AppName, "url": util.AppURL, "result": result}
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
