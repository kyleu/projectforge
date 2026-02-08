package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var screenResult = &Screen{Key: "result", Title: "Result Screen!", Hotkeys: []string{`"any key": back`, `"q": quit`}, Render: renderResult}

func renderResult(t *TUI) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Status"))
	b.WriteString("\n\n")

	body := t.result
	if body == "" {
		body = fmt.Sprintf("You selected:\n\n%s", t.choice)
	}

	result := resultStyle.Render(body)
	b.WriteString(result)

	content := b.String()

	return containerStyle.Width(t.width).Height(t.height).Render(content)
}

func onKeyReturn(scr *Screen) func(key string, t *TUI) tea.Cmd {
	return func(key string, t *TUI) tea.Cmd {
		if key == "q" {
			t.quitting = true
			return tea.Quit
		}
		t.Screen = scr
		return nil
	}
}
