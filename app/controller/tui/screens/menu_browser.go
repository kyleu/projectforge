package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
)

func clampMenuCursor(cursor int, count int) int {
	if count <= 0 {
		return 0
	}
	if cursor < 0 {
		return 0
	}
	if cursor >= count {
		return count - 1
	}
	return cursor
}

func moveMenuCursor(cursor int, count int, delta int) int {
	if count <= 0 || delta == 0 {
		return clampMenuCursor(cursor, count)
	}
	return clampMenuCursor(cursor+delta, count)
}

func menuMoveDelta(msg tea.Msg) (int, bool) {
	m, ok := msg.(tea.KeyMsg)
	if !ok {
		return 0, false
	}
	switch m.String() {
	case "up", "k":
		return -1, true
	case "down", "j":
		return 1, true
	default:
		return 0, false
	}
}

func panelContentSize(panelStyle lipgloss.Style, outerW int, outerH int) (int, int) {
	return ContentSize(panelStyle, outerW, outerH)
}

func renderMenuPanel(items menu.Items, cursor int, styles style.Styles, panelStyle lipgloss.Style, outerW int, outerH int) string {
	contentW, _ := panelContentSize(panelStyle, outerW, outerH)
	body := components.RenderMenuList(items, cursor, styles, contentW)
	return Bounded(panelStyle, outerW, outerH, body)
}

func renderTextPanel(body string, panelStyle lipgloss.Style, outerW int, outerH int) string {
	return Bounded(panelStyle, outerW, outerH, body)
}

func renderMainListScreen(title string, items menu.Items, cursor int, styles style.Styles, rects layout.Rects) string {
	contentW, _, _ := mainPanelContentSize(styles.Panel, rects)
	body := components.RenderMenuList(items, cursor, styles, contentW)
	return renderScreenPanel(title, body, styles.Panel, styles, rects)
}
