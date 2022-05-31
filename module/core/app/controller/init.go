// Package controller - $PF_IGNORE$
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

{{{ if.HasModule "marketing" }}}// Initialize system dependencies for the marketing site.
func initSite(_ *app.State, _ util.Logger) {
}

// Configure marketing site data for each request.
func initSiteRequest(_ *app.State, _ *cutil.PageState) error {
	return nil
}{{{ end }}}
