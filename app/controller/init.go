// Package controller $PF_IGNORE$
package controller

import (
	"path/filepath"
	"strings"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/util"
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
		for _, m := range prj.Modules {
			if strings.Contains(m, "@") {
				key, path := util.SplitStringLast(m, '@', true)
				destination := filepath.Join(prj.Path, path)
				_, added, err := as.Services.Modules.AddIfNeeded(key, destination)
				if err != nil {
					return errors.Wrapf(err, "unable to load referenced module [%s] from [%s]", key, destination)
				}
				if added {
					ps.Logger.Infof("added module [%s] using files in [%s]", key, destination)
				}
			}
		}
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
