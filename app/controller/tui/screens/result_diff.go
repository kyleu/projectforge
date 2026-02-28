package screens

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
)

const dataResultDiffScroll = "result_diff_scroll"

type ResultDiffScreen struct{}

func NewResultDiffScreen() *ResultDiffScreen {
	return &ResultDiffScreen{}
}

func (s *ResultDiffScreen) Key() string {
	return KeyResultDiff
}

func (s *ResultDiffScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	path := ps.EnsureData().GetStringOpt(dataResultDiffPath)
	if path == "" {
		ps.Title = "Diff"
		ps.SetStatus("No diff selected")
		return nil
	}
	ps.Title = fmt.Sprintf("Diff: %s", filepath.Base(path))
	ps.SetStatus("Viewing diff for [%s]", path)
	ps.EnsureData()[dataResultDiffScroll] = 0
	return nil
}

func (s *ResultDiffScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if delta, moved := menuMoveDelta(msg); moved {
		moveResultDiffScroll(ps, delta)
		return mvc.Stay(), nil, nil
	}

	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "pgdown", "J":
			moveResultDiffScroll(ps, 10)
			return mvc.Stay(), nil, nil
		case "pgup", "K":
			moveResultDiffScroll(ps, -10)
			return mvc.Stay(), nil, nil
		case "home", "g":
			ps.EnsureData()[dataResultDiffScroll] = 0
			return mvc.Stay(), nil, nil
		case "end", "G":
			lines := strings.Split(ps.EnsureData().GetStringOpt(dataResultDiffPatch), "\n")
			if len(lines) > 0 {
				ps.EnsureData()[dataResultDiffScroll] = len(lines) - 1
			}
			return mvc.Stay(), nil, nil
		case KeyEsc, KeyBackspace, KeyLeft, "h", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *ResultDiffScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	contentW, contentH, _ := mainPanelContentSize(styles.Panel, rects)
	bodyH := max(1, contentH-3)

	path := ps.EnsureData().GetStringOpt(dataResultDiffPath)
	tag := ps.EnsureData().GetStringOpt(dataResultDiffTag)
	patch := ps.EnsureData().GetStringOpt(dataResultDiffPatch)
	if patch == "" {
		patch = "(no patch available)"
	}

	header := fmt.Sprintf("%s %s", statusStyle(styles, statusKeyFromTag(tag)).Render(tag), path)
	divider := styles.Muted.Render(strings.Repeat("─", max(1, contentW-1)))
	content := s.renderPatchWindow(ps, bodyH, patch)
	body := strings.Join([]string{
		truncateLine(singleLine(header), max(1, contentW)),
		"Diff",
		divider,
		fitVertical(content, bodyH),
	}, "\n")

	return renderScreenPanel(ps.Title, body, styles.Panel, styles, rects)
}

func (s *ResultDiffScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: scroll", "pgup/pgdn: page", "g/G: top/bottom", "b: back"}}
}

func (s *ResultDiffScreen) renderPatchWindow(ps *mvc.PageState, height int, patch string) string {
	lines := strings.Split(patch, "\n")
	if len(lines) == 0 {
		return "(empty diff)"
	}
	scroll := ps.EnsureData().GetIntOpt(dataResultDiffScroll)
	scroll = max(0, min(scroll, len(lines)-1))
	ps.EnsureData()[dataResultDiffScroll] = scroll
	end := min(len(lines), scroll+max(1, height))
	return strings.Join(lines[scroll:end], "\n")
}

func moveResultDiffScroll(ps *mvc.PageState, delta int) {
	if delta == 0 {
		return
	}
	scroll := ps.EnsureData().GetIntOpt(dataResultDiffScroll)
	ps.EnsureData()[dataResultDiffScroll] = max(0, scroll+delta)
}

func statusKeyFromTag(tag string) string {
	switch tag {
	case "[A]":
		return "new"
	case "[D]":
		return "missing"
	case "[M]":
		return "different"
	default:
		return ""
	}
}
