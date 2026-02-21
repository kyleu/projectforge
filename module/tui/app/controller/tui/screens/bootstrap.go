// $PF_GENERATE_ONCE$
package screens

import (
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

func Bootstrap(_ *mvc.State) *Registry {
	reg := NewRegistry()

	reg.AddScreen(NewMainMenuScreen(reg)){{{ if .HasModule "docbrowse" }}}
	docsScreenItem := &menu.Item{Key: KeyDocs, Title: "Documentation", Description: "Browse embedded markdown documentation", Icon: "book", Route: KeyDocs}
	reg.Register(docsScreenItem, NewDocumentationScreen()){{{ end }}}
	settingsScreenItem := &menu.Item{Key: KeySettings, Title: "Settings", Description: "Runtime and diagnostics", Icon: "settings", Route: KeySettings}
	reg.Register(settingsScreenItem, NewSettingsScreen())
	aboutScreenItem := &menu.Item{Key: KeyAbout, Title: "About", Description: "Information about " + util.AppName, Icon: "info", Route: KeyAbout}
	reg.Register(aboutScreenItem, NewAboutScreen())

	// reg.AddScreen(SomeNewScreen())

	return reg
}
