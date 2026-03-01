package registry

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/screens/settings"
	"projectforge.dev/projectforge/app/controller/tui/screens/tproject"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

func Bootstrap(st *app.State, logger util.Logger) *screens.Registry {
	reg := screens.NewRegistry()

	reg.AddScreen(screens.NewMainMenuScreen(reg))

	const projectDesc = "Manage your projects"
	projectsScreenItem := &menu.Item{Key: tproject.KeyProjects, Title: "Projects", Description: projectDesc, Icon: "folder", Route: tproject.KeyProjects}
	reg.Register(projectsScreenItem, tproject.NewProjectsScreen())
	doctorScreenItem := &menu.Item{Key: screens.KeyDoctor, Title: "Doctor", Description: "Environment checks", Icon: "first-aid", Route: screens.KeyDoctor}
	reg.Register(doctorScreenItem, screens.NewDoctorScreen())

	docsScreenItem := &menu.Item{Key: screens.KeyDocs, Title: "Documentation", Description: "Browse markdown documentation", Icon: "book", Route: screens.KeyDocs}
	reg.Register(docsScreenItem, screens.NewDocumentationScreen())
	aboutScreenItem := &menu.Item{Key: screens.KeyAbout, Title: "About", Description: "Information about " + util.AppName, Icon: "info", Route: screens.KeyAbout}
	reg.Register(aboutScreenItem, screens.NewAboutScreen())

	reg.AddScreen(tproject.NewProjectScreen())
	reg.AddScreen(tproject.NewProjectNewScreen())
	reg.AddScreen(screens.NewFileBrowserScreen())
	reg.AddScreen(screens.NewFileViewerScreen())
	reg.AddScreen(tproject.NewResultsScreen())
	reg.AddScreen(tproject.NewResultDiffScreen())

	settings.Register(reg)

	return reg
}
