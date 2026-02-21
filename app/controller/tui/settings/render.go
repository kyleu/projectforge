package settings

import (
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/style"
)

func renderPanel(ts style.Styles, title string, body string, rects layout.Rects) string {
	contentW, contentH, outerH := panelDimensions(ts.Panel, rects)
	panel := ts.Panel.
		Width(contentW).
		Height(contentH).
		MaxWidth(max(1, rects.Main.W)).
		MaxHeight(outerH).
		Render(body)
	if !rects.Compact {
		return panel
	}
	header := ts.Header.Width(max(1, rects.Main.W)).Render(title)
	return lipgloss.JoinVertical(lipgloss.Left, header, panel)
}

func panelDimensions(panelStyle lipgloss.Style, rects layout.Rects) (int, int, int) {
	outerH := max(1, rects.Main.H)
	if rects.Compact {
		outerH = max(1, rects.Main.H-1)
	}
	contentW := max(1, rects.Main.W-panelStyle.GetHorizontalFrameSize())
	contentH := max(1, outerH-panelStyle.GetVerticalFrameSize())
	return contentW, contentH, outerH
}
