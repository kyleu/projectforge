package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
)

type ResultsScreen struct{}

func NewResultsScreen() *ResultsScreen {
	return &ResultsScreen{}
}

func (s *ResultsScreen) Key() string {
	return KeyResults
}

func (s *ResultsScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	prj := selectedProject(ts, ps)
	choice := actionChoiceFromData(ps.EnsureData())
	if choice.title == "" {
		ps.Title = "Results"
		ps.SetStatus("No action selected")
		return nil
	}
	ps.Title = fmt.Sprintf("%s Results", choice.title)
	ps.SetStatus("Running [%s]...", choice.title)
	message := ps.EnsureData().GetStringOpt(dataInputMessage)
	dryRun := ps.EnsureData().GetBoolOpt(dataInputDryRun)
	return runProjectActionCmd(ts, ps, prj, choice, message, dryRun)
}

func (s *ResultsScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	choice := actionChoiceFromData(ps.EnsureData())
	prj := selectedProject(ts, ps)
	message := ps.EnsureData().GetStringOpt(dataInputMessage)
	dryRun := ps.EnsureData().GetBoolOpt(dataInputDryRun)

	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "r", "enter":
			if choice.title == "" {
				return mvc.Stay(), nil, nil
			}
			ps.SetStatus("Running [%s]...", choice.title)
			return mvc.Stay(), runProjectActionCmd(ts, ps, prj, choice, message, dryRun), nil
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	case projectResultMsg:
		if m.err != nil {
			ps.SetError(m.err)
			return mvc.Stay(), nil, nil
		}
		ps.EnsureData()[dataResultLines] = m.lines
		ps.EnsureData()[dataResultTitle] = m.title
		ps.EnsureData()[dataResultComplete] = true
		ps.SetStatusText(m.title)
	}
	return mvc.Stay(), nil, nil
}

func (s *ResultsScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)

	lines := []string{"Running action..."}
	if ps.EnsureData().GetBoolOpt(dataResultComplete) {
		lines = ps.EnsureData().GetStringArrayOpt(dataResultLines)
		if len(lines) == 0 {
			lines = []string{"No output"}
		}
	}

	contentW, _, _ := mainPanelContentSize(styles.Panel, rects)
	body := renderLines(lines, contentW)
	return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
}

func (s *ResultsScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"r/enter: run again", "b: back"}}
}
