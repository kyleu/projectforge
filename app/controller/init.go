package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/pkg/errors"
)

// Initialize app-specific system dependencies.
func initApp() {
}

// Configure app-specific data for each request.
func initAppRequest(as *app.State, ps *cutil.PageState) error {
	_, err := as.Services.Projects.Refresh()
	if err != nil {
		return errors.Wrap(err, "can't load projects")
	}
	return nil
}

// Initialize system dependencies for the marketing site.
func initSite() {
}

// Configure marketing site data for each request.
func initSiteRequest(as *app.State, ps *cutil.PageState) error {
  return nil
}
