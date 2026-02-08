package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/util"
)

var screenSplash = NewScreen("splash", util.AppName, "", renderSplash, `"any key": continue`, `"q": quit`)

func renderSplash(t *TUI) string {
	var b strings.Builder

	b.WriteString(util.AppName)

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
