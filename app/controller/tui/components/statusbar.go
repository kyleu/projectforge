package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/style"
)

func RenderStatus(status string, err string, shortHelp []string, serverURL string, width int, st style.Styles) string {
	top := renderTopStatus(status, err, st)
	helpText := strings.Join(shortHelp, " | ")
	help := renderBottomStatus(helpText, serverURL, width, st)
	return top + "\n" + help
}

func renderTopStatus(status string, err string, st style.Styles) string {
	leftText := status
	leftStyle := st.Status.Padding(0)
	if err != "" {
		leftText = err
		leftStyle = st.Error.Padding(0)
	}
	leftText = strings.TrimSpace(leftText)
	if leftText == "" {
		leftText = " "
	}
	return leftStyle.Render(leftText)
}

func renderBottomStatus(helpText string, serverURL string, width int, st style.Styles) string {
	helpStyle := st.Muted.Padding(0)
	urlStyle := st.Muted.Padding(0).Underline(true)
	if width < 1 {
		if serverURL == "" {
			return helpStyle.Render(helpText)
		}
		return helpStyle.Render(helpText) + " " + urlStyle.Render(serverURL)
	}

	helpText = truncateStatusEllipsis(helpText, width)
	if serverURL == "" {
		return helpStyle.Render(helpText)
	}
	rightText := serverURL
	rightW := lipgloss.Width(rightText)
	if rightW >= width {
		return urlStyle.Render(truncateStatusEllipsis(rightText, width))
	}
	maxHelp := width - rightW - 1
	if maxHelp < 1 {
		return urlStyle.Render(truncateStatusEllipsis(rightText, width))
	}
	helpText = truncateStatusEllipsis(helpText, maxHelp)
	left := helpStyle.Render(helpText)
	right := urlStyle.Render(rightText)
	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
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
		return "…"
	}
	return string(r[:width-1]) + "…"
}
