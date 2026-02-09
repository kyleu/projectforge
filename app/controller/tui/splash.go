package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/util"
)

var screenSplash = NewScreen("splash", util.AppName, "", renderSplash, `"any key": continue`, `"q": quit`)

func renderSplash(t *TUI) string {
	var b strings.Builder

	art := util.ToASCIIArt(util.AppName)
	lines := strings.Split(strings.TrimRight(art, "\n"), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			b.WriteString("\n")
			continue
		}
		style := splashTitleStyle.Foreground(acct)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(splashSubtitleStyle.Render("Server is running at " + t.serverURL))
	b.WriteString("\n\n")
	b.WriteString(pressKeyStyle.Render("Press any key to continue..."))

	content := b.String()

	return containerStyle.Width(t.width).Height(t.height).Render(content)
}

func onKeySplash(key string, t *TUI) tea.Cmd {
	if key == "q" {
		return tea.Quit
	}
	t.Screen = screenMenu
	return nil
}
