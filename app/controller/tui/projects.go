package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
)

var screenProjects = NewScreen("projects", "Projects", "", renderProjects, keysMenu...)

type projectsLoadedMsg struct {
	err error
}

func renderProjects(t *TUI) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Projects"))
	b.WriteString("\n\n")

	prjs := projectsFor(t)
	switch {
	case t.Config.projectsLoading:
		b.WriteString("Loading projects...")
	case t.Config.projectsErr != nil:
		b.WriteString("Error loading projects:\n")
		b.WriteString(t.Config.projectsErr.Error())
	case len(prjs) == 0:
		b.WriteString("No projects found.")
	default:
		items := make(menu.Items, 0, len(prjs))
		for _, p := range prjs {
			items = append(items, &menu.Item{Key: p.Key, Title: p.Title()})
		}
		cursor := t.Screen.Cursor()
		b.WriteString(RenderMenuOptions(cursor, items))
		b.WriteString("\n")
		if cursor >= 0 && cursor < len(prjs) {
			b.WriteString(resultStyle.Render(projectDetails(prjs[cursor])))
		}
	}

	return containerStyle.Width(t.width).Height(t.height).Render(b.String())
}

func onKeyProjects(key string, t *TUI) tea.Cmd {
	switch key {
	case "q":
		t.quitting = true
		return tea.Quit
	case tuiKeyEsc:
		t.Screen = screenMenu
	case tuiKeyUp, "k":
		if t.Screen.Cursor() > 0 {
			t.Screen.ModifyCursor(-1)
		}
	case tuiKeyDown, "j":
		prjs := projectsFor(t)
		if t.Screen.Cursor() < len(prjs)-1 {
			t.Screen.ModifyCursor(1)
		}
	case "r":
		t.Config.projectsLoading = true
		t.Config.projectsErr = nil
		return loadProjectsCmd(t)
	}
	return nil
}

func loadProjectsCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		_, err := t.st.Services.Projects.Refresh(t.logger)
		return projectsLoadedMsg{err: err}
	}
}

func projectsFor(t *TUI) project.Projects {
	return t.st.Services.Projects.Projects()
}

func projectDetails(p *project.Project) string {
	if p == nil {
		return "No project selected."
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Name: %s\n", p.Title()))
	b.WriteString(fmt.Sprintf("Key: %s\n", p.Key))
	b.WriteString(fmt.Sprintf("Package: %s\n", orDash(p.Package)))
	b.WriteString(fmt.Sprintf("Version: %s\n", orDash(p.Version)))
	b.WriteString(fmt.Sprintf("Exec: %s\n", orDash(p.ExecSafe())))
	b.WriteString(fmt.Sprintf("Path: %s\n", orDash(p.Path)))
	return b.String()
}
