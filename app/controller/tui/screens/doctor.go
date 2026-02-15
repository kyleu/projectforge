package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

type DoctorScreen struct{}

type doctorResultMsg struct {
	title string
	lines []string
	err   error
}

func NewDoctorScreen() *DoctorScreen {
	return &DoctorScreen{}
}

func (s *DoctorScreen) Key() string {
	return KeyDoctor
}

func (s *DoctorScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Doctor"
	ps.SetStatus("Use 'a' to run all checks")
	if ts.App != nil && ts.App.Services != nil {
		checks.SetModules(ts.App.Services.Modules)
	}
	return nil
}

func (s *DoctorScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	ck := s.availableChecks(ts)
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			if ps.Cursor > 0 {
				ps.Cursor--
			}
		case "down", "j":
			if ps.Cursor < len(ck)-1 {
				ps.Cursor++
			}
		case "a":
			ps.SetStatus("Running all doctor checks...")
			return mvc.Stay(), s.runAll(ts, ps), nil
		case "r", "enter":
			if len(ck) == 0 {
				return mvc.Stay(), nil, nil
			}
			selected := ck[ps.Cursor]
			ps.SetStatus("Running check [%s]...", selected.Title)
			return mvc.Stay(), s.runOne(ps, selected), nil
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	case doctorResultMsg:
		if m.err != nil {
			ps.SetError(m.err)
			return mvc.Stay(), nil, nil
		}
		ps.EnsureData()["result"] = m.lines
		ps.SetStatusText(m.title)
	}
	return mvc.Stay(), nil, nil
}

func (s *DoctorScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	ck := s.availableChecks(ts)
	items := lo.Map(ck, func(c *doctor.Check, _ int) *menu.Item {
		return &menu.Item{Key: c.Key, Title: c.Title, Description: c.Summary}
	})

	title := styles.Header.Width(max(1, rects.Main.W)).Render(ps.Title)

	result := "Run checks to see detailed output"
	if lines := ps.EnsureData().GetStringArrayOpt("result"); len(lines) > 0 {
		result = strings.Join(lines, "\n")
	}

	bodyH := max(1, rects.Main.H-1)
	if rects.Compact {
		leftStyle := styles.Panel
		rightStyle := styles.Sidebar
		leftW := max(1, rects.Main.W-leftStyle.GetHorizontalFrameSize())
		rightW := max(1, rects.Main.W-rightStyle.GetHorizontalFrameSize())
		leftH := max(1, bodyH/2-leftStyle.GetVerticalFrameSize())
		rightH := max(1, bodyH-bodyH/2-rightStyle.GetVerticalFrameSize())
		list := leftStyle.Width(leftW).Height(leftH).Render(components.RenderMenuList(items, ps.Cursor, styles, leftW))
		output := rightStyle.Width(rightW).Height(rightH).Render(renderLines(strings.Split(result, "\n"), rightW))
		return lipgloss.JoinVertical(lipgloss.Left, title, list, output)
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
	list := leftStyle.Width(leftCW).Height(leftCH).Render(components.RenderMenuList(items, ps.Cursor, styles, leftCW))
	output := rightStyle.Width(rightCW).Height(rightCH).Render(renderLines(strings.Split(result, "\n"), rightCW))
	return lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.JoinHorizontal(lipgloss.Top, list, output))
}

func (s *DoctorScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "r/enter: run check", "a: run all", "b: back"}}
}

func (s *DoctorScreen) availableChecks(ts *mvc.State) doctor.Checks {
	prjs := ts.App.Services.Projects.Projects()
	checks.SetModules(ts.App.Services.Modules)
	return checks.ForModules(prjs.AllModules())
}

func (s *DoctorScreen) runAll(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	modules := ts.App.Services.Projects.Projects().AllModules()
	ctx := ps.Context
	logger := ps.Logger
	if logger == nil {
		logger = ts.Logger
	}
	return func() tea.Msg {
		results := checks.CheckAll(ctx, modules, logger)
		lines := flattenDoctorResults(results)
		return doctorResultMsg{title: fmt.Sprintf("Ran %d checks", len(results)), lines: lines}
	}
}

func (s *DoctorScreen) runOne(ps *mvc.PageState, c *doctor.Check) tea.Cmd {
	ctx := ps.Context
	logger := ps.Logger
	return func() tea.Msg {
		if c == nil {
			return doctorResultMsg{err: fmt.Errorf("check not found")}
		}
		ret := c.Check(ctx, logger)
		if ret == nil {
			return doctorResultMsg{err: fmt.Errorf("check [%s] not applicable", c.Key)}
		}
		lines := flattenDoctorResults(doctor.Results{ret})
		return doctorResultMsg{title: fmt.Sprintf("Finished [%s]", c.Title), lines: lines}
	}
}

func flattenDoctorResults(results doctor.Results) []string {
	lines := make([]string, 0, len(results)*3)
	for _, r := range results {
		if r == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("[%s] %s", r.Status, r.Title))
		for _, l := range r.Logs {
			lines = append(lines, "  - "+l)
		}
		for _, e := range r.Errors {
			lines = append(lines, "  - error: "+e.String())
		}
		for _, sol := range r.CleanSolutions() {
			lines = append(lines, "  - solution: "+sol)
		}
	}
	if len(lines) == 0 {
		return []string{util.OK}
	}
	return lines
}
