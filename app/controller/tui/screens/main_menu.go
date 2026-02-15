package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/util"
)

type MainMenuScreen struct {
	registry *Registry
}

func NewMainMenuScreen(registry *Registry) *MainMenuScreen {
	return &MainMenuScreen{registry: registry}
}

func (s *MainMenuScreen) Key() string {
	return KeyMainMenu
}

func (s *MainMenuScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = util.AppName
	ps.SetStatus("Choose a section")
	return nil
}

func (s *MainMenuScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	items := s.registry.Menu().Visible()
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if ps.Cursor > 0 {
				ps.Cursor--
			}
		case "down", "j":
			if ps.Cursor < len(items)-1 {
				ps.Cursor++
			}
		case "enter":
			if len(items) == 0 {
				return mvc.Stay(), nil, nil
			}
			item := items[ps.Cursor]
			return mvc.Push(item.Route, nil), nil, nil
		case "q":
			return mvc.Quit(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *MainMenuScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	title := styles.Header.Width(max(1, rects.Main.W)).Render(ps.Title)
	body := components.RenderMenuList(s.registry.Menu().Visible(), ps.Cursor, styles, max(1, rects.Main.W-4))
	panel := styles.Panel.Width(max(1, rects.Main.W)).Height(max(1, rects.Main.H-1)).Render(body)
	return lipgloss.JoinVertical(lipgloss.Left, title, panel)
}

func (s *MainMenuScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open", "q: quit"}}
}
