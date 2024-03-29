// Package controller - $PF_GENERATE_ONCE$
package controller

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

// Initialize system dependencies for the marketing site.
func initSite(context.Context, *app.State, util.Logger) error {
	return nil
}

// Configure marketing site data for each request.
func initSiteRequest(*app.State, *cutil.PageState) error {
	return nil
}
