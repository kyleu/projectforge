package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
)

type HelpBindings struct {
	Short []string
	Full  []string
}

type Screen interface {
	Key() string
	Init(*mvc.State, *mvc.PageState) tea.Cmd
	Update(*mvc.State, *mvc.PageState, tea.Msg) (mvc.Transition, tea.Cmd, error)
	View(*mvc.State, *mvc.PageState, layout.Rects) string
	Help(*mvc.State, *mvc.PageState) HelpBindings
}
