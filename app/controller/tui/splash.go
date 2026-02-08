package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/util"
)

var screenSplash = &Screen{Key: "splash", Title: util.AppName, Hotkeys: []string{`"any key": continue`, `"q": quit`}, Render: renderSplash}

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
