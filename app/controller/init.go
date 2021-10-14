// Package controller $PF_IGNORE$
package controller

import (
	"strings"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/project"
	"github.com/pkg/errors"
)

// Initialize app-specific system dependencies.
func initApp(*app.State) {
}

// Configure app-specific data for each request.
func initAppRequest(as *app.State, ps *cutil.PageState) error {
	prjs, err := as.Services.Projects.Refresh()
	if err != nil {
		return errors.Wrap(err, "can't load projects")
	}
	for _, prj := range prjs {
		var mods []*project.ModuleDef
		if prj.Info != nil {
			mods = prj.Info.ModuleDefs
		}
		keys, err := as.Services.Modules.Register(prj.Path, mods...)
		if err != nil {
			return errors.Wrap(err, "unable to register module definitions")
		}
		if len(keys) > 0 {
			as.Logger.Debugf("Loaded modules for [%s]: %s", prj.Key, strings.Join(keys, ", "))
		}
	}

	root := as.Services.Projects.ByPath(".")
	if root.Info == nil {
		ps.ForceRedirect = "/p/" + root.Key + "/edit"
	}

	return nil
}

// Initialize system dependencies for the marketing site.
func initSite(*app.State) {
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}
