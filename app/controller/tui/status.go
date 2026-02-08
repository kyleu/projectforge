package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/util"
)

var keysDefault = []string{`"q" quit`, `"↑"/"↓" move`, `"enter" select`}

func (t *TUI) withStatus(content string) string {
	return content + "\n" + t.renderStatusLine(t.width)
}

func (t *TUI) renderStatusLine(width int) string {
	versionPart := util.AppName + " v" + t.st.AppVersion()

	keysPart := util.StringJoin(util.Choose(len(t.Screen.Hotkeys) == 0, keysDefault, t.Screen.Hotkeys), " • ")
	leftText := fmt.Sprintf("%s • %s", versionPart, keysPart)
	rightText := t.serverURL

	contentWidth := width - 2 // statusLineStyle has Padding(0, 1)
	if contentWidth < 0 {
		contentWidth = 0
	}

	var text string
	if rightText == "" || contentWidth == 0 {
		text = truncate(leftText, contentWidth)
	} else {
		text = joinStatusParts(leftText, rightText, contentWidth)
	}
	return statusLineStyle.Width(width).Render(text)
}

func joinStatusParts(left string, right string, width int) string {
	if width <= 0 {
		return ""
	}

	rw := lipgloss.Width(right)
	if rw >= width {
		return truncate(right, width)
	}

	const sep = " "
	lw := width - rw - lipgloss.Width(sep)
	if lw <= 0 {
		return truncate(right, width)
	}

	left = truncate(left, lw)
	pad := lw - lipgloss.Width(left)
	if pad < 0 {
		pad = 0
	}
	return left + strings.Repeat(" ", pad) + sep + right
}

func truncate(s string, width int) string {
	if width <= 0 {
		return s
	}
	if len([]rune(s)) <= width {
		return s
	}
	rs := []rune(s)
	if width <= 1 {
		return string(rs[:width])
	}
	return string(rs[:width-1]) + "…"
}
