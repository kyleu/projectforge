// $PF_IGNORE$
package site

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/site/download"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/verror"
	"projectforge.dev/projectforge/views/vsite"
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
			"<meta name=\"go-import\" content=\"projectforge.dev/projectforge git %s\">\n  " +
			"<meta name=\"go-source\" content=\"projectforge.dev/projectforge %s %s/tree/master{/dir} %s/blob/master{/dir}/{file}#L{line}\">"
		ps.HeaderContent = fmt.Sprintf(msg, util.AppSource, util.AppSource, util.AppSource, util.AppSource)
		return "", &vsite.GoSource{}, path, nil
	case keyFeatures:
		switch {
		case len(path) == 1:
			page, err = featureList(as, ps)
		case len(path) == 2:
			page, err = featureDetail(path[1], as, ps)
		default:
			path, page, err = featureFiles(path, as, ps)
		}
	case keyAbout:
		ps.Title = "About " + util.AppName
		ps.Data = util.AppName + " v" + as.BuildInfo.Version
		page = &views.About{}
	case keyDownload:
		dls := download.GetLinks(as.BuildInfo.Version)
		ps.Title = "Downloads"
		ps.Data = util.ValueMap{"base": "https://github.com/kyleu/projectforge/releases/download/v" + as.BuildInfo.Version, "links": dls}
		page = &vsite.Download{Links: dls}
	case keyInstall:
		page, err = mdTemplate("Installation", "This static page contains installation instructions", "installation.md", "code", ps)
	case keyContrib:
		page, err = mdTemplate("Contributing", "This static page describes how to build "+util.AppName, "contributing.md", "cog", ps)
	case keyTech:
		page, err = mdTemplate("Technology", "This static page describes the technology used in "+util.AppName, "technology.md", "shield", ps)
	case keyFAQ:
		page, err = mdTemplate("FAQ", "Frequently asked questions about "+util.AppName, "faq.md", "question", ps)
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
