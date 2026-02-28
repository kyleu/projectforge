package screens

import (
	"time"

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

// SidebarContentProvider optionally lets a screen provide shell sidebar content.
// Return handled=false to hide the sidebar and use a full-width content pane.
type SidebarContentProvider interface {
	SidebarContent(*mvc.State, *mvc.PageState, layout.Rects) (content string, handled bool)
}

// AutoRefreshProvider optionally enables timed refresh ticks while the screen is active.
// Returning a non-positive duration disables auto refresh for the screen.
type AutoRefreshProvider interface {
	RefreshInterval(*mvc.State, *mvc.PageState) time.Duration
}

// GlobalKeyBypassProvider optionally allows screens to suppress root-level key handlers.
type GlobalKeyBypassProvider interface {
	BypassGlobalKey(*mvc.State, *mvc.PageState, string) bool
}
