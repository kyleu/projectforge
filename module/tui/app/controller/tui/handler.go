// $PF_GENERATE_ONCE$
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func handleMessage(t *TUI, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	default:
		return t, nil
	}
}
