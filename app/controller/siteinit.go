// Package controller - Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/help"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

// Initialize system dependencies for the marketing site.
func initSite(as *app.State, logger util.Logger) {
	mod, err := module.NewService(context.Background(), as.Files.Root(), logger)
	if err != nil {
		logger.Errorf("unable to initialize site: %+v", err)
	}
	hlp := help.NewService(logger)
	as.Services = &app.Services{Help: hlp, Modules: mod}
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}
