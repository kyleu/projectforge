package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type ProjectsScreen struct{}

const (
	keyNewProject       = "project:new"
	projectsSelectedKey = "projects.selectedKey"
)

func NewProjectsScreen() *ProjectsScreen {
	return &ProjectsScreen{}
}

func (s *ProjectsScreen) Key() string {
	return KeyProjects
}

func (s *ProjectsScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	prjs := ts.App.Services.Projects.Projects()
	items := projectItems(prjs)
	ps.Title = "Projects"
	if selected, ok := ps.EnsureData()[projectsSelectedKey].(string); ok && selected != "" {
		ps.Cursor = indexForProjectItem(items, selected)
	} else {
		ps.Cursor = util.Choose(len(prjs) > 0, 1, 0)
	}
	ps.EnsureData()[projectsSelectedKey] = items[clampMenuCursor(ps.Cursor, len(items))].Key
	ps.SetStatus("Loaded %d project(s)", len(prjs))
	return nil
}

func (s *ProjectsScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	items := projectItems(ts.App.Services.Projects.Projects())
	ps.Cursor = clampMenuCursor(ps.Cursor, len(items))
	ps.EnsureData()[projectsSelectedKey] = items[ps.Cursor].Key
	if delta, moved := menuMoveDelta(msg); moved {
		ps.Cursor = moveMenuCursor(ps.Cursor, len(items), delta)
		ps.EnsureData()[projectsSelectedKey] = items[ps.Cursor].Key
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok {
		switch m.String() {
		case KeyEnter:
			if len(items) == 0 {
				return mvc.Stay(), nil, nil
			}
			selected := items[ps.Cursor]
			ps.EnsureData()[projectsSelectedKey] = selected.Key
			if selected.Key == keyNewProject {
				return mvc.Push(KeyProjectNew, nil), nil, nil
			}
			return mvc.Push(KeyProject, util.ValueMap{"project": selected.Key}), nil, nil
		case KeyEsc, KeyBackspace, "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectsScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	prjs := ts.App.Services.Projects.Projects()
	items := projectItems(prjs)
	return renderMainListScreen(ps.Title, items, clampMenuCursor(ps.Cursor, len(items)), styles, rects)
}

func (s *ProjectsScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open project/create", "b: back"}}
}

func projectItems(prjs project.Projects) menu.Items {
	ret := make(menu.Items, 0, len(prjs)+1)
	ret = append(ret, &menu.Item{Key: keyNewProject, Title: "[+] New Project", Description: "Guided project creation form"})
	for _, prj := range prjs {
		desc := strings.TrimSpace(prj.DescriptionSafe())
		if desc == "" {
			desc = fmt.Sprintf("%d module(s)", len(prj.Modules))
		}
		ret = append(ret, &menu.Item{Key: prj.Key, Title: prj.Title(), Description: desc})
	}
	return ret
}

func indexForProjectItem(items menu.Items, key string) int {
	for idx, item := range items {
		if item.Key == key {
			return idx
		}
	}
	return 0
}
