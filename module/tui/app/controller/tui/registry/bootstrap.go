// $PF_GENERATE_ONCE$
package registry

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/screens/settings"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

func Bootstrap(st *app.State, logger util.Logger) *screens.Registry {
	reg := screens.NewRegistry()

	reg.AddScreen(screens.NewMainMenuScreen(reg)){{{ if .HasModule "docbrowse" }}}

	docsScreenItem := &menu.Item{Key: screens.KeyDocs, Title: "Documentation", Description: "Browse markdown documentation", Icon: "book", Route: screens.KeyDocs}
	reg.Register(docsScreenItem, screens.NewDocumentationScreen()){{{ end }}}
	aboutScreenItem := &menu.Item{Key: screens.KeyAbout, Title: "About", Description: "Information about " + util.AppName, Icon: "info", Route: screens.KeyAbout}
	reg.Register(aboutScreenItem, screens.NewAboutScreen())

	// reg.AddScreen(SomeNewScreen())

	settings.Register(reg)

	return reg
}
