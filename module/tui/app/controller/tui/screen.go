package tui

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/util"
)

type Screen struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Help    string   `json:"help,omitzero"`
	Hotkeys []string `json:"keys,omitzero"`

	Update func(msg tea.Msg) tea.Cmd        `json:"-"`
	Render func(t *TUI) string              `json:"-"`
	OnKey  func(key string, t *TUI) tea.Cmd `json:"-"`

	State util.ValueMap `json:"state,omitzero"`
}

func NewScreen(key string, title string, help string, render func(t *TUI) string, hotkeys ...string) *Screen {
	return &Screen{Key: key, Title: title, Help: help, Render: render, Hotkeys: hotkeys, State: util.ValueMap{}}
}

func (s *Screen) Cursor() int {
	return s.State.GetIntOpt("cursor")
}

func (s *Screen) ModifyCursor(delta int) {
	s.State["cursor"] = s.State.GetIntOpt("cursor") + delta
}

func (s *Screen) ResetCursor() {
	s.State["cursor"] = 0
}

var screensOnce sync.Once

func initScreensIfNeeded() {
	screensOnce.Do(func() {
		// $PF_SECTION_START(tui-init)$
		screenSplash.OnKey = onKeySplash
		screenResult.OnKey = onKeyReturn(screenMenu)
		// $PF_SECTION_END(tui-init)$
	})
}
