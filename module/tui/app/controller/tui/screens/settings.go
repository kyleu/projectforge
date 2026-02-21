package screens

import (
	"fmt"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/style"
)

type SettingsScreen struct{}

func NewSettingsScreen() *SettingsScreen {
	return &SettingsScreen{}
}

func (s *SettingsScreen) Key() string {
	return KeySettings
}

func (s *SettingsScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Settings"
	ps.SetStatus("Runtime diagnostics")
	if ts.App != nil {
		ps.EnsureData()["settings"] = []string{
			fmt.Sprintf("debug: %t", ts.App.Debug),
			fmt.Sprintf("goos/goarch: %s/%s", runtime.GOOS, runtime.GOARCH),
			fmt.Sprintf("server url: %s", ts.ServerURL),
		}
	}
	return nil
}

func (s *SettingsScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok {
		switch m.String() {
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *SettingsScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	lines := ps.EnsureData().GetStringArrayOpt("settings")
	if len(lines) == 0 {
		lines = []string{"No settings available"}
	}
	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
}

func (s *SettingsScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"b: back"}}
}
