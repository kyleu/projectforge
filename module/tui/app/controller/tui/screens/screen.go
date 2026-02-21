package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
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

// SidebarContentProvider optionally lets a screen override the shell sidebar content.
// Return handled=false to fall back to the default sidebar.
type SidebarContentProvider interface {
	SidebarContent(*mvc.State, *mvc.PageState, layout.Rects) (content string, handled bool)
}
