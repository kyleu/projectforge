package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/lib/menu"
)

var screenProjects = NewScreen("projects", "Projects", "", renderProjects, `"esc": back`, `"↑"/"↓": move`, `"r": reload`, `"q": quit`)

func renderProjects(t *TUI) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Projects"))
	b.WriteString("\n\n")

	switch {
	case t.projectsLoading:
		b.WriteString("Loading projects...")
	case t.projectsErr != nil:
		b.WriteString("Error loading projects:\n")
		b.WriteString(t.projectsErr.Error())
	case len(t.projects) == 0:
		b.WriteString("No projects found.")
	default:
		items := make(menu.Items, 0, len(t.projects))
		for _, p := range t.projects {
			items = append(items, &menu.Item{Key: p.Key, Title: p.Title()})
		}
		cursor := t.Screen.Cursor()
		b.WriteString(RenderMenuOptions(cursor, items))
		b.WriteString("\n")
		if cursor >= 0 && cursor < len(t.projects) {
			b.WriteString(resultStyle.Render(projectDetails(t.projects[cursor])))
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
		if t.Screen.Cursor() < len(t.projects)-1 {
			t.Screen.ModifyCursor(1)
		}
	case "r":
		t.projectsLoading = true
		t.projectsErr = nil
		return loadProjectsCmd(t)
	}
	return nil
}
