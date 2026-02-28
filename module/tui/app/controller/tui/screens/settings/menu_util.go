package settings

import (
	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/components"
	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/menu"
)

func menuClamp(cursor int, count int) int {
	return screens.MenuClampCursor(cursor, count)
}

func menuDelta(msg tea.Msg) (int, bool) {
	return screens.MenuMoveDelta(msg)
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
