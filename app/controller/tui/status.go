package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/util"
)

var keysDefault = []string{`"↑"/"↓" move`, `"enter" select`, `"/" logs`, `"q" quit`}

func (t *TUI) withStatus(content string) string {
	lines := []string{content}
	if t.showLogs {
		if logLines := t.logLineLimit(); logLines > 0 {
			lines = append(lines, t.renderLogPanel(t.width, logLines))
		}
	}
	lines = append(lines, t.renderStatusLine(t.width))
	return strings.Join(lines, "\n")
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

func (t *TUI) renderLogPanel(width int, limit int) string {
	var b strings.Builder
	b.WriteString(renderDivider(width))
	b.WriteByte('\n')
	b.WriteString(logHeaderStyle.Width(width).Render("Recent logs"))
	b.WriteByte('\n')

	logs := t.latestLogs(limit)
	if len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}

	for i := 0; i < limit; i++ {
		if i >= len(logs) {
			b.WriteString(logLineStyle.Width(width).Render(""))
			b.WriteByte('\n')
			continue
		}
		log := logs[i]
		line := strings.Join(strings.Fields(log), " ")
		b.WriteString(logLineStyle.Width(width).Render(truncate(line, width)))
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderDivider(width int) string {
	if width <= 0 {
		return ""
	}
	return logDividerStyle.Width(width).Render(strings.Repeat("─", width))
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
