package tui

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

type TUI struct {
	Config *Config
	Screen *Screen

	choice string
	result string

	width    int
	height   int
	showLogs bool

	logs   []string
	logsMu sync.RWMutex

	ctx       context.Context //nolint:containedctx // lifecycle-owned by TUI; used to cancel background work
	logger    util.Logger
	st        *app.State
	serverURL string

	quitting bool
}

func NewTUI(ctx context.Context, st *app.State, serverURL string, logger util.Logger) *TUI {
	initScreensIfNeeded()
	return &TUI{Config: &Config{}, Screen: screenSplash, ctx: ctx, logger: logger, st: st, serverURL: serverURL}
}

func (t *TUI) Init() tea.Cmd {
	return nil
}

func (t *TUI) Run(ctx context.Context, logger util.Logger) error {
	logger.Debugf("Starting TUI for [%s]", util.AppName)
	p := tea.NewProgram(t, tea.WithAltScreen(), tea.WithContext(ctx))
	_, err := p.Run()
	return err
}

func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
		return t, nil
	default:
		return handleMessage(t, msg)
	}
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

	reservedLines := 1 // status bar
	if t.showLogs {
		reservedLines += t.logPanelHeight()
	}
	contentHeight := t.height - reservedLines
	if contentHeight < 1 {
		contentHeight = 1
	}

	originalHeight := t.height
	t.height = contentHeight
	content := t.Screen.Render(t)
	t.height = originalHeight

	return t.withStatus(content)
}

func (t *TUI) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	if key == tuiKeyCtrlC {
		t.quitting = true
		return t, tea.Quit
	}
	if key == "/" {
		t.showLogs = !t.showLogs
		return t, nil
	}
	if t.Screen == nil {
		t.quitting = true
		return t, tea.Quit
	}
	return t, t.Screen.OnKey(key, t)
}

func (t *TUI) AddLog(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) {
	t.logsMu.Lock()
	defer t.logsMu.Unlock()

	logMsg := message // for now

	t.logs = append(t.logs, logMsg)
	if len(t.logs) > 50 {
		t.logs = t.logs[1:]
	}
}

func (t *TUI) latestLogs(limit int) []string {
	if limit <= 0 {
		return nil
	}

	t.logsMu.RLock()
	defer t.logsMu.RUnlock()

	if len(t.logs) == 0 {
		return nil
	}

	start := len(t.logs) - limit
	if start < 0 {
		start = 0
	}

	ret := make([]string, len(t.logs)-start)
	copy(ret, t.logs[start:])
	return ret
}

func (t *TUI) logPanelHeight() int {
	lines := t.logLineLimit()
	if !t.showLogs || lines == 0 {
		return 0
	}
	return lines + 2 // divider + header + log lines
}

func (t *TUI) logLineLimit() int {
	if !t.showLogs {
		return 0
	}
	maxLogLines := 5
	available := t.height - 2 // leave at least one line for main content and one for status
	if available <= 2 {
		return 0
	}
	lines := available - 2 // divider + header
	if lines > maxLogLines {
		lines = maxLogLines
	}
	return lines
}
