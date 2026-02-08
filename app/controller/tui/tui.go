package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type TUI struct {
	Screen *Screen

	choice string
	result string

	width  int
	height int

	ctx       context.Context //nolint:containedctx // lifecycle-owned by TUI; used to cancel background work
	logger    util.Logger
	st        *app.State
	serverURL string

	quitting bool
}

func NewTUI(ctx context.Context, st *app.State, serverURL string, logger util.Logger) *TUI {
	initScreensIfNeeded()

	return &TUI{
		Screen:    screenSplash,
		ctx:       ctx,
		logger:    logger,
		st:        st,
		serverURL: serverURL,
	}
}

func (t *TUI) Init() tea.Cmd {
	return nil
}

func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t *TUI) View() string {
	if t.quitting {
		return ""
	}
	if t.width == 0 {
		t.width = 80
	}
	if t.height == 0 {
		t.height = 24
	}
	if t.height == 0 {
		t.height = 1
	}
	return t.withStatus(t.Screen.Render(t))
}

func (t *TUI) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == tuiKeyCtrlC {
		t.quitting = true
		return t, tea.Quit
	}
	if t.Screen == nil {
		t.quitting = true
		return t, tea.Quit
	}
	return t, t.Screen.OnKey(key, t)
}
