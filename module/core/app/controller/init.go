// Package controller - $PF_GENERATE_ONCE$
package controller

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

// Initialize app-specific system dependencies.
func initApp(_ *app.State, _ util.Logger) {
}

// Configure app-specific data for each request.
func initAppRequest(_ *app.State, _ *cutil.PageState) error {
	return nil
}
