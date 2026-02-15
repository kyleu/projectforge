package tui

import (
	"context"
	"sync"
	"time"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

type TUI struct {
	st        *app.State
	serverURL string

	logger util.Logger

	logsMu sync.RWMutex
	logs   []string
}

func NewTUI(st *app.State, serverURL string, logger util.Logger) *TUI {
	return &TUI{logger: logger, st: st, serverURL: serverURL}
}

func (t *TUI) Run(ctx context.Context, logger util.Logger) error {
	logger.Debugf("starting [%s] TUI...", util.AppName)
	return nil
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
