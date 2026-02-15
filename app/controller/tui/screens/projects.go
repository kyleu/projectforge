package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type ProjectsScreen struct{}

func NewProjectsScreen() *ProjectsScreen {
	return &ProjectsScreen{}
}

func (s *ProjectsScreen) Key() string {
	return KeyProjects
}

func (s *ProjectsScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Projects"
	ps.SetStatus("Loaded %d project(s)", len(ts.App.Services.Projects.Projects()))
	return nil
}

func (s *ProjectsScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	prjs := ts.App.Services.Projects.Projects()
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if ps.Cursor > 0 {
				ps.Cursor--
			}
		case "down", "j":
			if ps.Cursor < len(prjs)-1 {
				ps.Cursor++
			}
		case "enter":
			if len(prjs) == 0 {
				return mvc.Stay(), nil, nil
			}
			prj := prjs[ps.Cursor]
			return mvc.Push(KeyProject, util.ValueMap{"project": prj.Key}), nil, nil
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectsScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	prjs := ts.App.Services.Projects.Projects()
	title := styles.Header.Width(max(1, rects.Main.W)).Render(ps.Title)
	list := components.RenderMenuList(projectItems(prjs), ps.Cursor, styles, max(1, rects.Main.W-4))
	panel := styles.Panel.Width(max(1, rects.Main.W)).Height(max(1, rects.Main.H-1)).Render(list)
	return lipgloss.JoinVertical(lipgloss.Left, title, panel)
}

func (s *ProjectsScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open project", "b: back"}}
}

func projectItems(prjs project.Projects) menu.Items {
	ret := make(menu.Items, 0, len(prjs))
	for _, prj := range prjs {
		desc := strings.TrimSpace(prj.DescriptionSafe())
		if desc == "" {
			desc = fmt.Sprintf("%d module(s)", len(prj.Modules))
		}
		ret = append(ret, &menu.Item{Key: prj.Key, Title: prj.Title(), Description: desc})
	}
	return ret
}
