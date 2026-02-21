package screens

import (
	"github.com/charmbracelet/lipgloss"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/style"
)

func mainPanelOuterHeight(rects layout.Rects) int {
	if rects.Compact {
		return max(1, rects.Main.H-1)
	}
	return max(1, rects.Main.H)
}

func mainPanelContentSize(panelStyle lipgloss.Style, rects layout.Rects) (int, int, int) {
	outerH := mainPanelOuterHeight(rects)
	contentW, contentH := panelContentSize(panelStyle, rects.Main.W, outerH)
	return contentW, contentH, outerH
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

func ContentSize(st lipgloss.Style, outerW int, outerH int) (int, int) {
	w := max(1, outerW-st.GetHorizontalFrameSize())
	h := max(1, outerH-st.GetVerticalFrameSize())
	return w, h
}

func Bounded(st lipgloss.Style, outerW int, outerH int, content string) string {
	ow := max(1, outerW)
	oh := max(1, outerH)
	cw, ch := ContentSize(st, ow, oh)
	return st.
		Width(cw).
		Height(ch).
		MaxWidth(ow).
		MaxHeight(oh).
		Render(content)
}
