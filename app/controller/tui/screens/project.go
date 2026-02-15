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
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

type ProjectScreen struct{}

type projectResultMsg struct {
	title string
	lines []string
	err   error
}

type actionChoice struct {
	key         string
	title       string
	description string
}

func NewProjectScreen() *ProjectScreen {
	return &ProjectScreen{}
}

func (s *ProjectScreen) Key() string {
	return KeyProject
}

func (s *ProjectScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	prj := selectedProject(ts, ps)
	if prj == nil {
		ps.Title = "Project"
		ps.SetStatus("Project not found")
		return nil
	}
	ps.Title = prj.Title()
	ps.SetStatus("Choose an action for [%s]", prj.Key)
	return nil
}

func (s *ProjectScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	prj := selectedProject(ts, ps)
	choices := s.choices()
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if ps.Cursor > 0 {
				ps.Cursor--
			}
		case "down", "j":
			if ps.Cursor < len(choices)-1 {
				ps.Cursor++
			}
		case "enter", "r":
			if prj == nil || len(choices) == 0 {
				return mvc.Stay(), nil, nil
			}
			choice := choices[ps.Cursor]
			ps.SetStatus("Running [%s]...", choice.title)
			return mvc.Stay(), s.runAction(ts, ps, prj, choice), nil
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	case projectResultMsg:
		if m.err != nil {
			ps.SetError(m.err)
			return mvc.Stay(), nil, nil
		}
		ps.EnsureData()["result"] = m.lines
		ps.SetStatusText(m.title)
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	prj := selectedProject(ts, ps)
	title := ps.Title
	if prj == nil {
		title = "Project Not Found"
	}
	header := styles.Header.Width(max(1, rects.Main.W)).Render(title)
	items := make(menu.Items, 0, len(s.choices()))
	for _, c := range s.choices() {
		items = append(items, &menu.Item{Key: c.key, Title: c.title, Description: c.description})
	}

	metaLines := []string{"No project loaded"}
	if prj != nil {
		metaLines = []string{
			fmt.Sprintf("Project: %s", prj.Title()),
			fmt.Sprintf("Path: %s", prj.Path),
		}
		metaLines = append(metaLines, moduleLines(prj.Modules, 8)...)
	}
	if lines := ps.EnsureData().GetStringArrayOpt("result"); len(lines) > 0 {
		metaLines = append(metaLines, "", "Result:")
		metaLines = append(metaLines, lines...)
	}

	bodyH := max(1, rects.Main.H-1)
	if rects.Compact {
		leftStyle := styles.Panel
		rightStyle := styles.Sidebar
		leftW := max(1, rects.Main.W-leftStyle.GetHorizontalFrameSize())
		rightW := max(1, rects.Main.W-rightStyle.GetHorizontalFrameSize())
		leftH := max(1, bodyH/2-leftStyle.GetVerticalFrameSize())
		rightH := max(1, bodyH-bodyH/2-rightStyle.GetVerticalFrameSize())
		left := leftStyle.Width(leftW).Height(leftH).Render(components.RenderMenuList(items, ps.Cursor, styles, leftW))
		right := rightStyle.Width(rightW).Height(rightH).Render(renderLines(metaLines, rightW))
		return lipgloss.JoinVertical(lipgloss.Left, header, left, right)
	}

	leftW := max(24, (rects.Main.W*2)/3)
	if leftW > rects.Main.W-20 {
		leftW = max(1, rects.Main.W-20)
	}
	rightW := max(1, rects.Main.W-leftW)
	leftStyle := styles.Panel
	rightStyle := styles.Sidebar
	leftCW := max(1, leftW-leftStyle.GetHorizontalFrameSize())
	rightCW := max(1, rightW-rightStyle.GetHorizontalFrameSize())
	leftCH := max(1, bodyH-leftStyle.GetVerticalFrameSize())
	rightCH := max(1, bodyH-rightStyle.GetVerticalFrameSize())
	left := leftStyle.Width(leftCW).Height(leftCH).Render(components.RenderMenuList(items, ps.Cursor, styles, leftCW))
	right := rightStyle.Width(rightCW).Height(rightCH).Render(renderLines(metaLines, rightCW))
	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, left, right))
}

func (s *ProjectScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: run", "b: back"}}
}

func (s *ProjectScreen) choices() []actionChoice {
	ret := make([]actionChoice, 0, 16)
	for _, t := range action.ProjectTypes {
		ret = append(ret, actionChoice{key: "action:" + t.Key, title: t.Title, description: t.Description})
	}
	for _, t := range []action.Type{action.TypeDebug, action.TypeRules, action.TypeSVG} {
		ret = append(ret, actionChoice{key: "action:" + t.Key, title: t.Title, description: t.Description})
	}
	for _, ga := range []*git.Action{git.ActionStatus, git.ActionFetch, git.ActionPull, git.ActionPush, git.ActionCommit, git.ActionReset, git.ActionHistory, git.ActionMagic} {
		ret = append(ret, actionChoice{key: "git:" + ga.Key, title: ga.Title, description: ga.Description})
	}
	return ret
}

func (s *ProjectScreen) runAction(ts *mvc.State, ps *mvc.PageState, prj *project.Project, c actionChoice) tea.Cmd {
	ctx := ps.Context
	logger := ps.Logger
	if logger == nil {
		logger = ts.Logger
	}
	return func() tea.Msg {
		if prj == nil {
			return projectResultMsg{err: fmt.Errorf("project not found")}
		}
		if strings.HasPrefix(c.key, "action:") {
			act := action.TypeFromString(strings.TrimPrefix(c.key, "action:"))
			params := &action.Params{
				ProjectKey: prj.Key,
				T:          act,
				Cfg:        util.ValueMap{},
				MSvc:       ts.App.Services.Modules,
				PSvc:       ts.App.Services.Projects,
				ESvc:       ts.App.Services.Export,
				XSvc:       ts.App.Services.Exec,
				Logger:     logger,
			}
			res := action.Apply(ctx, params)
			lines := append([]string{fmt.Sprintf("status: %s", res.Status)}, res.Logs...)
			lines = append(lines, res.Errors...)
			return projectResultMsg{title: fmt.Sprintf("Completed [%s]", act.Title), lines: lines}
		}
		gs := git.NewService(prj.Key, prj.Path)
		key := strings.TrimPrefix(c.key, "git:")
		var (
			res *git.Result
			err error
		)
		switch key {
		case git.ActionStatus.Key:
			res, err = gs.Status(ctx, logger)
		case git.ActionFetch.Key:
			res, err = gs.Fetch(ctx, logger)
		case git.ActionPull.Key:
			res, err = gs.Pull(ctx, logger)
		case git.ActionPush.Key:
			res, err = gs.Push(ctx, logger)
		case git.ActionCommit.Key:
			res, err = gs.Commit(ctx, "Project Forge TUI commit", logger)
		case git.ActionReset.Key:
			res, err = gs.Reset(ctx, logger)
		case git.ActionHistory.Key:
			res, err = gs.History(ctx, &git.HistoryArgs{Limit: 25}, logger)
		case git.ActionMagic.Key:
			res, err = gs.Magic(ctx, "Project Forge TUI magic", true, logger)
		default:
			err = fmt.Errorf("unknown git action [%s]", key)
		}
		if err != nil {
			return projectResultMsg{err: err}
		}
		if res == nil {
			return projectResultMsg{err: fmt.Errorf("no result returned")}
		}
		lines := []string{fmt.Sprintf("status: %s", res.Status)}
		for _, k := range res.CleanData().Keys() {
			lines = append(lines, fmt.Sprintf("%s: %v", k, res.CleanData()[k]))
		}
		if res.Error != "" {
			lines = append(lines, "error: "+res.Error)
		}
		return projectResultMsg{title: fmt.Sprintf("Completed [%s]", c.title), lines: lines}
	}
}

func selectedProject(ts *mvc.State, ps *mvc.PageState) *project.Project {
	if ts == nil || ts.App == nil || ts.App.Services == nil || ts.App.Services.Projects == nil {
		return nil
	}
	key := ps.EnsureData().GetStringOpt("project")
	if key == "" {
		return ts.App.Services.Projects.Default()
	}
	prj, err := ts.App.Services.Projects.Get(key)
	if err != nil {
		return nil
	}
	return prj
}

func moduleLines(mods []string, maxShown int) []string {
	if len(mods) == 0 {
		return []string{"Modules: (none)"}
	}
	if maxShown < 1 {
		maxShown = 1
	}
	ret := []string{"Modules:"}
	limit := len(mods)
	if limit > maxShown {
		limit = maxShown
	}
	for _, mod := range mods[:limit] {
		ret = append(ret, "  - "+mod)
	}
	if len(mods) > maxShown {
		ret = append(ret, fmt.Sprintf("  - ... (+%d more)", len(mods)-maxShown))
	}
	return ret
}

func renderLines(lines []string, width int) string {
	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		ret = append(ret, truncateLine(singleLine(line), width))
	}
	return strings.Join(ret, "\n")
}

func truncateLine(s string, width int) string {
	if width < 1 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
