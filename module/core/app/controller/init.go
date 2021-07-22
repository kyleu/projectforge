// Package controller $PF_IGNORE$
package controller

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller/cutil"
)

// Initialize app-specific system dependencies.
func initApp(*app.State) {
}

// Configure app-specific data for each request.
func initAppRequest(*app.State, *cutil.PageState) error {
	return nil
}

// Initialize system dependencies for the marketing site.
func initSite(*app.State) {
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}
