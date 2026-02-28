// $PF_HAS_MODULE(task)$
package settings

import (
	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

const (
	keyTaskDetail = "settings.admin.task.detail"
	keyTaskRun    = "settings.admin.task.run"
)

type taskListScreen struct{ items menu.Items }

func registerTaskScreens(reg *screens.Registry) {
	reg.AddScreen(&taskListScreen{})
	reg.AddScreen(&taskDetailScreen{})
	reg.AddScreen(&taskRunScreen{})
}

func (s *taskListScreen) Key() string { return keyTask }

func (s *taskListScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Task Engine"
	s.items = nil
	for _, t := range ts.App.Services.Task.RegisteredTasks {
		s.items = append(s.items, &menu.Item{Key: t.Key, Title: t.TitleSafe(), Description: t.Description, Route: keyTaskDetail})
	}
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	ps.SetStatusText(util.StringPlural(len(s.items), "Task"))
	return nil
}

func (s *taskListScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if delta, ok := menuDelta(msg); ok {
		ps.Cursor = menuClamp(ps.Cursor+delta, len(s.items))
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == screens.KeyEnter && len(s.items) > 0 {
		item := s.items[menuClamp(ps.Cursor, len(s.items))]
		return mvc.Push(keyTaskDetail, util.ValueMap{"key": item.Key}), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == screens.KeyEsc || m.String() == screens.KeyBackspace || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *taskListScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	return renderPanel(st, ps.Title, renderMenuBody(s.items, ps.Cursor, st, rects), rects)
}

func (s *taskListScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"up/down: move", "enter: detail", "b: back"}}
}
