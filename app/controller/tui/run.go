package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func Run(ctx context.Context, st *app.State, serverURL string, logger util.Logger) error {
	logger.Debugf("Starting TUI for [%s]", util.AppName)

	t := NewTUI(ctx, st, serverURL, logger)
	p := tea.NewProgram(t, tea.WithAltScreen(), tea.WithContext(ctx))

	_, err := p.Run()
	return err
}
