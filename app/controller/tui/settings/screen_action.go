package settings

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
)

type actionFn func(*mvc.State, *mvc.PageState) ([]string, error)

type actionScreen struct {
	key   string
	title string
	run   actionFn
}

func newActionScreen(key string, title string, run actionFn) *actionScreen {
	return &actionScreen{key: key, title: title, run: run}
}
func (s *actionScreen) Key() string { return s.key }

func (s *actionScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = s.title
	lines, err := s.run(ts, ps)
	if err != nil {
		return func() tea.Msg { return errMsg{err: err} }
	}
	ps.EnsureData()["lines"] = lines
	ps.SetStatus("Action completed")
	return nil
}

func (s *actionScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *actionScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	return renderPanel(st, ps.Title, strings.Join(ps.EnsureData().GetStringArrayOpt("lines"), "\n"), rects)
}

func (s *actionScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"b: back"}}
}
