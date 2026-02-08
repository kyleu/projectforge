package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

var (
	keysMenu   = []string{`"esc": back`, `"↑"/"↓" move`, `"enter" select`, `"q" quit`}
	screenMenu = NewScreen("menu", "Menu Screen!", "", renderMenu, keysMenu...)
)

var MainMenuItems = menu.Items{
	// $PF_SECTION_START(tui-menu)$
	// $PF_SECTION_END(tui-menu)$
}

func renderMenu(t *TUI) string {
	var b strings.Builder

	header := titleStyle.Render(util.AppName)
	b.WriteString(header)
	b.WriteString("\n\n")

	b.WriteString("Select an option:\n\n")

	b.WriteString(RenderMenuOptions(t.Screen.Cursor(), MainMenuItems))

	content := b.String()

	return containerStyle.Width(t.width).Height(t.height).Render(content)
}

func RenderMenuOptions(cursor int, items menu.Items) string {
	width := menuOptionWidth(items)
	itemStyle := menuItemStyle.Width(width)
	selectedStyle := menuSelectedStyle.Width(width)

	var b strings.Builder
	for i, it := range items {
		cursorStr := tuiCursorSpacer
		style := itemStyle
		if i == cursor {
			cursorStr = menuCursorStyle.Render("→ ")
			style = selectedStyle
		}

		b.WriteString(cursorStr)
		b.WriteString(style.Render(it.Title))
		b.WriteByte('\n')
	}

	return b.String()
}

func menuOptionWidth(items menu.Items) int {
	return stringWidth(items.Titles())
}

func stringWidth(items []string) int {
	maxLabelWidth := 0
	for _, it := range items {
		w := lipgloss.Width(it)
		if w > maxLabelWidth {
			maxLabelWidth = w
		}
	}

	frameWidth := menuItemStyle.GetHorizontalFrameSize()
	selectedFrameWidth := menuSelectedStyle.GetHorizontalFrameSize()
	if selectedFrameWidth > frameWidth {
		frameWidth = selectedFrameWidth
	}

	return maxLabelWidth + frameWidth + 2
}

func onKeyMenu(key string, t *TUI) tea.Cmd {
	cursor := t.Screen.Cursor()
	switch key {
	case "q":
		t.quitting = true
		return tea.Quit
	case "up", "k":
		if cursor > 0 {
			t.Screen.ModifyCursor(-1)
		}
	case "down", "j":
		if cursor < len(MainMenuItems)-1 {
			t.Screen.ModifyCursor(1)
		}
	case "enter", " ":
		t.choice = MainMenuItems[cursor].Title

		switch MainMenuItems[cursor].Key {
		case "projects":
			t.Screen = screenProjects
			t.Screen.ResetCursor()
			t.projectsLoading = true
			t.projectsErr = nil
			return loadProjectsCmd(t)
		case "doctor":
			t.Screen = screenDoctor
			t.Screen.ResetCursor()
			t.doctorLoading = true
			t.doctorRunning = false
			t.doctorErr = nil
			return loadDoctorChecksCmd(t)
		case "quit":
			t.quitting = true
			return tea.Quit
		default:
			t.Screen = screenResult
			t.result = fmt.Sprintf("You selected:\n\n%s", t.choice)
		}
	}
	return nil
}
