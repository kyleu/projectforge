// Package controller - $PF_GENERATE_ONCE$
package controller

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

// Initialize system dependencies for the marketing site.
func initSite(*app.State, util.Logger) {
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}
