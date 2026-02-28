package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/util"
)

func RenderStatus(status string, err string, shortHelp []string, serverURL string, serverErr string, width int, st style.Styles) string {
	top := renderTopStatus(status, err, st)
	helpText := strings.Join(shortHelp, " | ")
	help := renderBottomStatus(helpText, serverURL, serverErr, width, st)
	return top + "\n" + help
}

func renderTopStatus(status string, err string, st style.Styles) string {
	leftText := strings.TrimSpace(status)
	leftStyle := st.Status.Padding(0)
	if err != "" {
		leftText = strings.TrimSpace(err)
		leftStyle = st.Error.Padding(0)
	}
	if leftText == "" {
		leftText = " "
	}
	return leftStyle.Render(leftText)
}

func renderBottomStatus(helpText string, serverURL string, serverErr string, width int, st style.Styles) string {
	helpStyle := st.Muted.Padding(0)
	rightStyle := st.Muted.Padding(0).Underline(true)
	rightText := serverURL
	if serverErr != "" {
		rightStyle = st.Error.Padding(0)
		rightText = lastErrorSegment(serverErr)
	}
	if width < 1 {
		if rightText == "" {
			return helpStyle.Render(helpText)
		}
		return helpStyle.Render(helpText) + helpStyle.Render(" ") + rightStyle.Render(rightText)
	}

	helpText = truncateStatusEllipsis(helpText, width)
	if rightText == "" {
		return helpStyle.Render(helpText)
	}
	rightW := lipgloss.Width(rightText)
	if rightW >= width {
		return rightStyle.Render(truncateStatusEllipsis(rightText, width))
	}
	maxHelp := width - rightW - 1
	if maxHelp < 1 {
		return rightStyle.Render(truncateStatusEllipsis(rightText, width))
	}
	helpText = truncateStatusEllipsis(helpText, maxHelp)
	left := helpStyle.Render(helpText)
	right := rightStyle.Render(rightText)
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + helpStyle.Render(strings.Repeat(" ", gap)) + right
}

func lastErrorSegment(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx < 0 {
		return strings.TrimSpace(s)
	}
	ret := strings.TrimSpace(s[idx+1:])
	if ret == "" {
		return strings.TrimSpace(s)
	}
	return ret
}

func truncateStatusEllipsis(s string, width int) string {
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
