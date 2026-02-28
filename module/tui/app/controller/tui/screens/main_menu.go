package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

const mainMenuExitKey = "mainmenu.exit"

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
	items := s.items()
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
			if item.Key == mainMenuExitKey {
				return mvc.Quit(), nil, nil
			}
			return mvc.Push(item.Route, nil), nil, nil
		case "q", "esc":
			return mvc.Quit(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *MainMenuScreen) SidebarContent(ts *mvc.State, ps *mvc.PageState, _ layout.Rects) (string, bool) {
	styles := style.New(ts.Theme)
	items := s.items()
	cursor := clampMenuCursor(ps.Cursor, len(items))

	lines := []string{fmt.Sprintf("%s TUI", util.AppName), ""}
	lines = AppendSidebarProp(lines, styles, "sections", len(items))
	if len(items) > 0 {
		item := items[cursor]
		lines = AppendSidebarProp(lines, styles, "selected", item.Title)
		lines = AppendSidebarProp(lines, styles, "route", item.Route)
		if item.Description != "" {
			lines = AppendSidebarProp(lines, styles, "about", item.Description)
		}
	}
	lines = append(lines, "", "keys:", "up/down move", "enter open", "esc/q quit")
	return strings.Join(lines, "\n"), true
}

func (s *MainMenuScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	items := s.items()
	return renderMainListScreen(ps.Title, items, clampMenuCursor(ps.Cursor, len(items)), styles, rects)
}

func (s *MainMenuScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open", "esc|q: quit"}}
}

func (s *MainMenuScreen) items() menu.Items {
	items := s.registry.Menu().Visible()
	ret := make(menu.Items, 0, len(items)+1)
	ret = append(ret, items...)
	ret = append(ret, &menu.Item{Key: mainMenuExitKey, Title: "Exit", Description: "Exit the TUI"})
	return ret
}
