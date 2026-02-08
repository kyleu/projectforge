package tui

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Help    string   `json:"help,omitzero"`
	Hotkeys []string `json:"keys,omitzero"`

	Render func(t *TUI) string              `json:"-"`
	OnKey  func(key string, t *TUI) tea.Cmd `json:"-"`

	Cursor int
}

var screensOnce sync.Once

func initScreensIfNeeded() {
	screensOnce.Do(func() {
		screenMenu.OnKey = onKeyMenu

		screenSplash.OnKey = onKeySplash
		screenResult.OnKey = onKeyReturn(screenMenu)
	})
}
