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
		page, err = mdTemplate("This static page contains installation instructions", "installation.md", "code", ps)
	case keyContrib:
		page, err = mdTemplate("This static page describes how to build "+util.AppName, "contributing.md", "cog", ps)
	case keyTech:
		page, err = mdTemplate("This static page describes the technology used in "+util.AppName, "technology.md", "shield", ps)
	default:
		page, err = mdTemplate("Documentation for "+util.AppName, path[0]+".md", "", ps)
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

func mdTemplate(description string, path string, icon string, ps *cutil.PageState) (layout.Page, error) {
	if icon == "" {
		icon = "cog"
	}
	title := strings.TrimSuffix(path, ".md")
	html, err := doc.HTML(path, func(s string) (string, error) {
		ret, err := cutil.FormatMarkdown(s)
		if err != nil {
			return "", err
		}
		if h1Idx := strings.Index(ret, "<h1>"); h1Idx > -1 {
			if h1EndIdx := strings.Index(ret, "</h1>"); h1EndIdx > -1 {
				title = s[h1Idx+4 : h1EndIdx]
			}
			ic := fmt.Sprintf(`<svg class="icon" style="width: 20px; height: 20px;"><use xlink:href="#svg-%s" /></svg> `, icon)
			ret = ret[:h1Idx+4] + ic + ret[h1Idx+4:]
			ret = strings.ReplaceAll(ret, "<h3>", "<h4>")
			ret = strings.ReplaceAll(ret, "</h3>", "</h4>")
			ret = strings.ReplaceAll(ret, "<h2>", "<h4>")
			ret = strings.ReplaceAll(ret, "</h2>", "</h4>")
			ret = strings.ReplaceAll(ret, "<h1>", "<h3 style=\"margin-top: 0;\">")
			ret = strings.ReplaceAll(ret, "</h1>", "</h3>")
		}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	ps.Data = siteData(title, "description", description)
	ps.Title = title
	page := &vsite.MarkdownPage{Title: title, HTML: html}
	return page, nil
}
