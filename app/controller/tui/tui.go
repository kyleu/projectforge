package tui

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/tui/framework"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/util"
)

type TUI struct {
	st        *app.State
	serverURL string

	logger util.Logger

	logsMu sync.RWMutex
	logs   []string
}

func NewTUI(st *app.State, serverURL string, logger util.Logger) (*TUI, error) {
	return InitTUI(&TUI{logger: logger, st: st, serverURL: serverURL})
}

func (t *TUI) Run(ctx context.Context, logger util.Logger) error {
	logger.Debugf("starting [%s] TUI...", util.AppName)
	ts := mvc.NewState(ctx, t.st, t.serverURL, logger)
	registry := screens.Bootstrap(ts)
	root, err := framework.NewRootModel(ts, registry, screens.KeyMainMenu)
	if err != nil {
		return err
	}
	program := tea.NewProgram(root, tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = program.Run()
	return err
}

func (t *TUI) AddLog(_ string, _ time.Time, _ string, message string, _ util.ValueMap, _ string, _ util.ValueMap) {
	t.logsMu.Lock()
	defer t.logsMu.Unlock()
	t.logs = append(t.logs, message)
	if len(t.logs) > 200 {
		t.logs = t.logs[1:]
	}
}
