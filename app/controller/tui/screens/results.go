package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/util"
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
	changes := resultChanges(ps)
	complete := ps.EnsureData().GetBoolOpt(dataResultComplete)

	if complete {
		ps.Cursor = clampMenuCursor(ps.Cursor, len(changes))
	}

	if delta, moved := menuMoveDelta(msg); moved && complete && len(changes) > 0 {
		ps.Cursor = moveMenuCursor(ps.Cursor, len(changes), delta)
		return mvc.Stay(), nil, nil
	}

	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "r":
			if choice.title == "" {
				return mvc.Stay(), nil, nil
			}
			ps.SetStatus("Running [%s]...", choice.title)
			return mvc.Stay(), runProjectActionCmd(ts, ps, prj, choice, message, dryRun), nil
		case "enter":
			if complete && len(changes) > 0 {
				ch := changes[clampMenuCursor(ps.Cursor, len(changes))]
				return mvc.Push(KeyResultDiff, resultDiffData(ch)), nil, nil
			}
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
			return mvc.Stay(), nil, m.err
		}
		ps.EnsureData()[dataResultLines] = m.lines
		ps.EnsureData()[dataResultChanges] = m.changes
		ps.EnsureData()[dataResultStatus] = m.status
		ps.EnsureData()[dataResultTitle] = m.title
		ps.EnsureData()[dataResultComplete] = true
		ps.Cursor = 0
		ps.SetStatusText(m.title)
	}
	return mvc.Stay(), nil, nil
}

func (s *ResultsScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	contentW, _, _ := mainPanelContentSize(styles.Panel, rects)

	lines := []string{"Running action..."}
	if ps.EnsureData().GetBoolOpt(dataResultComplete) {
		changes := resultChanges(ps)
		if len(changes) > 0 {
			body := s.renderChanges(styles, ps.EnsureData().GetStringOpt(dataResultStatus), ps.Cursor, changes, contentW)
			return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
		}
		lines = ps.EnsureData().GetStringArrayOpt(dataResultLines)
		if len(lines) == 0 {
			lines = []string{"No output"}
		}
	}

	body := renderLines(lines, contentW)
	return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
}

func (s *ResultsScreen) Help(_ *mvc.State, ps *mvc.PageState) HelpBindings {
	if ps.EnsureData().GetBoolOpt(dataResultComplete) && len(resultChanges(ps)) > 0 {
		return HelpBindings{Short: []string{"up/down: select file", "enter: view diff", "r: run again", "b: back"}}
	}
	return HelpBindings{Short: []string{"r/enter: run again", "b: back"}}
}

func (s *ResultsScreen) renderChanges(styles style.Styles, status string, cursor int, changes []resultChange, width int) string {
	lines := make([]string, 0, len(changes)+3)
	lines = append(
		lines,
		"Status: "+util.OrDefault(status, "ok"),
		fmt.Sprintf("Changes: %s", util.StringPlural(len(changes), "change")),
		styles.Muted.Render(strings.Repeat("─", max(1, width-1))),
	)

	for i, ch := range changes {
		tag := "[" + resultChangeTag(ch.StatusKey) + "]"
		pathWidth := max(1, width-runeLen(tag)-1)
		row := fmt.Sprintf("%s %s", tag, truncateLine(ch.Path, pathWidth))
		if i == cursor {
			lines = append(lines, styles.Selected.Render(row))
			continue
		}
		lines = append(lines, statusStyle(styles, ch.StatusKey).Render(tag)+" "+truncateLine(ch.Path, pathWidth))
	}
	return strings.Join(lines, "\n")
}

func resultChanges(ps *mvc.PageState) []resultChange {
	raw, ok := ps.EnsureData()[dataResultChanges]
	if !ok || raw == nil {
		return nil
	}
	if ret, ok := raw.([]resultChange); ok {
		return ret
	}
	return nil
}

func resultDiffData(ch resultChange) util.ValueMap {
	return util.ValueMap{
		dataResultDiffPath:  ch.Path,
		dataResultDiffTag:   "[" + resultChangeTag(ch.StatusKey) + "]",
		dataResultDiffPatch: ch.Patch,
	}
}

func statusStyle(styles style.Styles, statusKey string) lipgloss.Style {
	switch statusKey {
	case diff.StatusNew.Key:
		return styles.Status.Bold(true)
	case diff.StatusMissing.Key:
		return styles.Error
	case diff.StatusDifferent.Key:
		return styles.App.Foreground(lipgloss.Color("#d7af00")).Bold(true)
	default:
		return styles.Muted.Bold(true)
	}
}
