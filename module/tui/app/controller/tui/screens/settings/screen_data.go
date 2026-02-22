package settings

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
)

type linesLoader func(*mvc.State, *mvc.PageState) ([]string, error)

type dataScreen struct {
	key    string
	title  string
	loader linesLoader
}

func newDataScreen(key string, title string, loader linesLoader) *dataScreen {
	return &dataScreen{key: key, title: title, loader: loader}
}

func (s *dataScreen) Key() string { return s.key }

func (s *dataScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = s.title
	lines, err := s.loader(ts, ps)
	if err != nil {
		return func() tea.Msg { return errMsg{err: err} }
	}
	ps.EnsureData()["lines"] = lines
	ps.SetStatus("Loaded")
	return nil
}

func (s *dataScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *dataScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	w, _, _ := panelDimensions(st.Panel, rects)
	raw := ps.EnsureData().GetStringArrayOpt("lines")
	lines := make([]string, len(raw))
	copy(lines, raw)
	if w > 0 {
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, w, "…")
		}
	}
	body := strings.Join(lines, "\n")
	return renderPanel(st, ps.Title, body, rects)
}

func (s *dataScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"b: back"}}
}
