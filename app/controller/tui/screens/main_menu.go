package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

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
	ps.Cursor = clampMenuCursor(ps.Cursor, len(items))
	if delta, moved := menuMoveDelta(msg); moved {
		ps.Cursor = moveMenuCursor(ps.Cursor, len(items), delta)
		return mvc.Stay(), nil, nil
	}
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
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

func (s *MainMenuScreen) SidebarContent(_ *mvc.State, ps *mvc.PageState, _ layout.Rects) (string, bool) {
	items := s.registry.Menu().Visible()
	cursor := clampMenuCursor(ps.Cursor, len(items))

	lines := []string{
		fmt.Sprintf("%s TUI", util.AppName),
		"",
		fmt.Sprintf("sections: %d", len(items)),
	}
	if len(items) > 0 {
		item := items[cursor]
		lines = append(lines,
			fmt.Sprintf("selected: %s", item.Title),
			fmt.Sprintf("route: %s", item.Route),
		)
		if item.Description != "" {
			lines = append(lines, fmt.Sprintf("about: %s", item.Description))
		}
	}
	lines = append(lines, "", "keys:", "up/down move", "enter open", "q quit")
	return strings.Join(lines, "\n"), true
}

func (s *MainMenuScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	items := s.registry.Menu().Visible()
	return renderMainListScreen(ps.Title, items, clampMenuCursor(ps.Cursor, len(items)), styles, rects)
}

func (s *MainMenuScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open", "q: quit"}}
}
