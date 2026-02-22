package settings

import (
	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
)

type menuScreen struct{ items menu.Items }

func newMenuScreen() *menuScreen {
	items := make(menu.Items, 0, len(adminOptions))
	for _, opt := range adminOptions {
		items = append(items, &menu.Item{Key: opt.Key, Title: opt.Title, Description: opt.Description, Route: opt.Key})
	}
	return &menuScreen{items: items}
}

func (s *menuScreen) Key() string { return screens.KeySettings }

func (s *menuScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Settings"
	ps.SetStatus("Choose an admin option")
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	return nil
}

func (s *menuScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	if delta, ok := menuDelta(msg); ok {
		ps.Cursor = menuClamp(ps.Cursor+delta, len(s.items))
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "enter" && len(s.items) > 0 {
		return mvc.Push(s.items[ps.Cursor].Route, nil), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *menuScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	body := renderMenuBody(s.items, ps.Cursor, st, rects)
	return renderPanel(st, ps.Title, body, rects)
}

func (s *menuScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"up/down: move", "enter: open", "b: back"}}
}
