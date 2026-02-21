package screens

import (
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/lib/menu"
)

func Bootstrap(_ *mvc.State) *Registry {
	reg := NewRegistry()

	reg.AddScreen(NewMainMenuScreen(reg))
	reg.Register(&menu.Item{Key: KeyProjects, Title: "Projects", Description: "Manage project workspaces", Icon: "folder", Route: KeyProjects}, NewProjectsScreen())
	reg.Register(&menu.Item{Key: KeyDocs, Title: "Documentation", Description: "Browse embedded markdown documentation", Icon: "book", Route: KeyDocs}, NewDocumentationScreen())
	reg.Register(&menu.Item{Key: KeyDoctor, Title: "Doctor", Description: "Dependency and environment checks", Icon: "first-aid", Route: KeyDoctor}, NewDoctorScreen())
	reg.Register(&menu.Item{Key: KeySettings, Title: "Settings", Description: "Runtime and diagnostics", Icon: "settings", Route: KeySettings}, NewSettingsScreen())
	reg.Register(&menu.Item{Key: KeyAbout, Title: "About", Description: "Build and environment metadata", Icon: "info", Route: KeyAbout}, NewAboutScreen())

	reg.AddScreen(NewProjectScreen())
	reg.AddScreen(NewResultsScreen())
	return reg
}
