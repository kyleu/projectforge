package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

const (
	keyExecNew    = "settings.admin.exec.new"
	keyExecDetail = "settings.admin.exec.detail"
)

type execListScreen struct{ items menu.Items }

func registerExecScreens(reg *screens.Registry) {
	reg.AddScreen(&execListScreen{})
	reg.AddScreen(&execNewScreen{})
	reg.AddScreen(&execDetailScreen{})
}

func (s *execListScreen) Key() string { return keyExec }

func (s *execListScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Managed Processes"
	s.items = menu.Items{{Key: keyExecNew, Title: "Start New Process", Description: "Open a form to create a new process", Route: keyExecNew}}
	for _, ex := range ts.App.Services.Exec.Execs {
		desc := fmt.Sprintf("pid=%d running=%t", ex.PID, ex.Completed == nil)
		s.items = append(s.items, &menu.Item{Key: ex.String(), Title: ex.String(), Description: desc, Route: keyExecDetail})
	}
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	ps.SetStatusText(util.StringPlural(len(ts.App.Services.Exec.Execs), "Process"))
	return nil
}

func (s *execListScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if delta, ok := menuDelta(msg); ok {
		ps.Cursor = menuClamp(ps.Cursor+delta, len(s.items))
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "enter" {
		item := s.items[menuClamp(ps.Cursor, len(s.items))]
		if item.Route == keyExecDetail {
			return mvc.Push(keyExecDetail, util.ValueMap{"key": item.Title}), nil, nil
		}
		return mvc.Push(keyExecNew, nil), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *execListScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	return renderPanel(st, ps.Title, renderMenuBody(s.items, ps.Cursor, st, rects), rects)
}

func (s *execListScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"up/down: move", "enter: open", "b: back"}}
}
