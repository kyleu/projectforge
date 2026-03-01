package tproject

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

type HelpBindings = screens.HelpBindings

const (
	KeyProjects   = "projects"
	KeyProject    = "project"
	KeyProjectNew = "project.new"
	KeyResults    = "results"
	KeyResultDiff = "result.diff"
)

const (
	dataFileExplorerTitle = "file_explorer_title"
	dataFileExplorerRoot  = "file_explorer_root"
	dataFileExplorerStart = "file_explorer_start"
)

func fileBrowserData(title string, root string, start string) util.ValueMap {
	return util.ValueMap{
		dataFileExplorerTitle: util.OrDefault(strings.TrimSpace(title), "Files"),
		dataFileExplorerRoot:  strings.TrimSpace(root),
		dataFileExplorerStart: strings.TrimSpace(start),
	}
}

func clampMenuCursor(cursor int, count int) int {
	return screens.MenuClampCursor(cursor, count)
}

func moveMenuCursor(cursor int, count int, delta int) int {
	return screens.MenuMoveCursor(cursor, count, delta)
}

func menuMoveDelta(msg tea.Msg) (int, bool) {
	return screens.MenuMoveDelta(msg)
}

func renderLines(lines []string, width int) string {
	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		ret = append(ret, truncateLine(singleLine(line), width))
	}
	return strings.Join(ret, "\n")
}

func truncateLine(s string, width int) string {
	if width < 1 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return util.KeyEllipsis
	}
	return string(r[:width-1]) + util.KeyEllipsis
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func runeLen(s string) int {
	return len([]rune(s))
}

func mainPanelOuterHeight(rects layout.Rects) int {
	if rects.Compact {
		return max(1, rects.Main.H-1)
	}
	return max(1, rects.Main.H)
}

func mainPanelContentSize(panelStyle lipgloss.Style, rects layout.Rects) (int, int, int) {
	outerH := mainPanelOuterHeight(rects)
	contentW, contentH := screens.ContentSize(panelStyle, rects.Main.W, outerH)
	return contentW, contentH, outerH
}

func renderTextPanel(body string, panelStyle lipgloss.Style, outerW int, outerH int) string {
	return screens.Bounded(panelStyle, outerW, outerH, body)
}

func renderScreenPanel(title string, body string, panelStyle lipgloss.Style, styles style.Styles, rects layout.Rects) string {
	_, _, outerH := mainPanelContentSize(panelStyle, rects)
	panel := renderTextPanel(body, panelStyle, rects.Main.W, outerH)
	return renderScreenFrame(title, panel, styles, rects)
}

func renderScreenFrame(title string, content string, styles style.Styles, rects layout.Rects) string {
	if !rects.Compact {
		return content
	}
	header := styles.Header.Width(max(1, rects.Main.W)).Render(title)
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func renderMainListScreen(title string, items menu.Items, cursor int, styles style.Styles, rects layout.Rects) string {
	contentW, _, _ := mainPanelContentSize(styles.Panel, rects)
	body := components.RenderMenuList(items, cursor, styles, contentW)
	return renderScreenPanel(title, body, styles.Panel, styles, rects)
}

func fitVertical(s string, height int) string {
	if height <= 1 {
		if s == "" {
			return ""
		}
		if i := strings.IndexByte(s, '\n'); i >= 0 {
			return s[:i]
		}
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) >= height {
		return strings.Join(lines[:height], "\n")
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func assertScreens() {
	var _ screens.Screen = NewProjectsScreen()
	var _ screens.Screen = NewProjectScreen()
	var _ screens.Screen = NewProjectNewScreen()
	var _ screens.Screen = NewResultsScreen()
	var _ screens.Screen = NewResultDiffScreen()
	_ = fmt.Sprintf
}
