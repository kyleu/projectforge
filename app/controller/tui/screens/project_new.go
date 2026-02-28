package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
)

type ProjectNewScreen struct{}

func NewProjectNewScreen() *ProjectNewScreen {
	return &ProjectNewScreen{}
}

func (s *ProjectNewScreen) Key() string {
	return KeyProjectNew
}

func (s *ProjectNewScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "New Project"
	return nil
}

func (s *ProjectNewScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectNewScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	return "TODO"
}

func (s *ProjectNewScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{}
}
