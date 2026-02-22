package tui

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/tui/framework"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/util"
)

type TUI struct {
	st        *app.State
	serverURL string
	serverErr string

	logger util.Logger

	registry *screens.Registry

	logsMu sync.RWMutex
	logs   []string
}

func NewTUI(st *app.State, serverURL string, serverErr string, registry *screens.Registry, logger util.Logger) (*TUI, error) {
	return InitTUI(&TUI{logger: logger, st: st, serverURL: serverURL, serverErr: serverErr, registry: registry})
}

func (t *TUI) Run(ctx context.Context, logger util.Logger) error {
	logger.Debugf("starting [%s] TUI...", util.AppName)
	if t.registry == nil {
		return errors.New("screen registry must be provided")
	}
	ts := mvc.NewState(ctx, t.st, t.serverURL, t.serverErr, logger, t.LastLogs)
	root, err := framework.NewRootModel(ts, t.registry, screens.KeyMainMenu)
	if err != nil {
		return err
	}
	program := tea.NewProgram(root, tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = program.Run()
	return err
}
