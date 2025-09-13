package cproject

import (
	"net/http"
	"net/url"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vpage"
)

const dblpipe = "||"

func GetProject(r *http.Request, as *app.State) (*project.Project, error) {
	key, err := cutil.PathString(r, "key", true)
	if err != nil {
		return nil, err
	}

	prj, err := as.Services.Projects.Get(key)
	if err != nil {
		return nil, err
	}
	return prj, nil
}

func GetProjectWithArgs(r *http.Request, as *app.State, logger util.Logger) (*project.Project, error) {
	prj, err := GetProject(r, as)
	if err != nil {
		return nil, err
	}
	err = prj.ModuleArgExport(as.Services.Projects, logger)
	if err != nil {
		return nil, err
	}
	return prj, nil
}

func HandleLoad(cfg util.ValueMap, u *url.URL, title string) layout.Page {
	if cfg.GetStringOpt("hasloaded") == util.BoolTrue {
		return nil
	}
	cutil.URLAddQuery(u, "hasloaded", util.BoolTrue)
	page := &vpage.Load{URL: u.String(), Title: title}
	return page
}
