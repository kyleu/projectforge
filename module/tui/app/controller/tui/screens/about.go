package screens

import (
	"fmt"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/util"
)

type AboutScreen struct{}

func NewAboutScreen() *AboutScreen {
	return &AboutScreen{}
}

func (s *AboutScreen) Key() string {
	return KeyAbout
}

func (s *AboutScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "About"
	started := ""
	if ts.App != nil {
		started = ts.App.Started.Format(time.RFC3339)
	}
	ps.EnsureData()["about"] = []string{
		fmt.Sprintf("name: %s", util.AppName),
		fmt.Sprintf("version: %s", ts.App.AppVersion()),
		fmt.Sprintf("summary: %s", util.AppSummary),
		fmt.Sprintf("url: %s", util.AppURL),
		fmt.Sprintf("source: %s", util.AppSource),
		fmt.Sprintf("runtime: %s/%s", runtime.GOOS, runtime.GOARCH),
		fmt.Sprintf("started: %s", started),
	}
	ps.SetStatus("Application metadata")
	return nil
}

func (s *AboutScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok {
		switch m.String() {
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *AboutScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	lines := ps.EnsureData().GetStringArrayOpt("about")
	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
}

func (s *AboutScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"b: back"}}
}
