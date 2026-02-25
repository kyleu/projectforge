package settings

import (
	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
)

func menuClamp(cursor int, count int) int {
	if count <= 0 || cursor < 0 {
		return 0
	}
	if cursor >= count {
		return count - 1
	}
	return cursor
}

func menuDelta(msg tea.Msg) (int, bool) {
	if m, ok := msg.(tea.MouseMsg); ok && m.Action == tea.MouseActionPress {
		switch m.Button {
		case tea.MouseButtonWheelUp:
			return -1, true
		case tea.MouseButtonWheelDown:
			return 1, true
		}
	}
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

func menuWindow(items menu.Items, cursor int, height int) (menu.Items, int) {
	if len(items) == 0 || height <= 0 || len(items) <= height {
		return items, menuClamp(cursor, len(items))
	}
	cursor = menuClamp(cursor, len(items))
	start := cursor - (height / 2)
	if start < 0 {
		start = 0
	}
	maxStart := len(items) - height
	if start > maxStart {
		start = maxStart
	}
	end := start + height
	return items[start:end], cursor - start
}

func renderMenuBody(items menu.Items, cursor int, st style.Styles, rects layout.Rects) string {
	w, h, _ := panelDimensions(st.Panel, rects)
	winItems, winCursor := menuWindow(items, cursor, h)
	return components.RenderMenuList(winItems, winCursor, st, w)
}
