package screens

import (
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

func Bootstrap(_ *mvc.State) *Registry {
	reg := NewRegistry()

	reg.AddScreen(NewMainMenuScreen(reg))
	projectsScreenItem := &menu.Item{Key: KeyProjects, Title: "Projects", Description: "Manage project workspaces", Icon: "folder", Route: KeyProjects}
	reg.Register(projectsScreenItem, NewProjectsScreen())
	docsScreenItem := &menu.Item{Key: KeyDocs, Title: "Documentation", Description: "Browse embedded markdown documentation", Icon: "book", Route: KeyDocs}
	reg.Register(docsScreenItem, NewDocumentationScreen())
	doctorScreenItem := &menu.Item{Key: KeyDoctor, Title: "Doctor", Description: "Dependency and environment checks", Icon: "first-aid", Route: KeyDoctor}
	reg.Register(doctorScreenItem, NewDoctorScreen())
	aboutScreenItem := &menu.Item{Key: KeyAbout, Title: "About", Description: "Information about " + util.AppName, Icon: "info", Route: KeyAbout}
	reg.Register(aboutScreenItem, NewAboutScreen())

	reg.AddScreen(NewProjectScreen())
	reg.AddScreen(NewResultsScreen())

	return reg
}
