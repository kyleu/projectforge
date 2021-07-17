package controller

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller/cutil"
)

// Initialize app-specific system dependencies.
func initApp() {
}

// Configure app-specific data for each request.
func initAppRequest(as *app.State, ps *cutil.PageState) error {
	return nil
}

// Initialize system dependencies for the marketing site.
func initSite() {
}

// Configure marketing site data for each request.
func initSiteRequest(as *app.State, ps *cutil.PageState) error {
	return nil
}
