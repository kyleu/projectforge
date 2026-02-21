package framework

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/util"
)

type navEntry struct {
	screen screens.Screen
	page   *mvc.PageState
}

type RootModel struct {
	state    *mvc.State
	registry *screens.Registry
	stack    []navEntry
	width    int
	height   int
	styles   style.Styles
	showLogs bool
}

const logDrawerLines = 8

func NewRootModel(state *mvc.State, registry *screens.Registry, initialScreen string) (*RootModel, error) {
	s, err := registry.Screen(initialScreen)
	if err != nil {
		return nil, err
	}
	ps := mvc.NewPageState(state.Context, initialScreen, strings.Title(initialScreen), nil, state.Logger)
	ret := &RootModel{
		state:    state,
		registry: registry,
		stack:    []navEntry{{screen: s, page: ps}},
		styles:   style.New(state.Theme),
	}
	return ret, nil
}

func (m *RootModel) Init() tea.Cmd {
	curr := m.current()
	if curr.screen == nil {
		return nil
	}
	return curr.screen.Init(m.state, curr.page)
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = x.Width, x.Height
		return m, tea.ClearScreen
	case tea.KeyMsg:
		if x.String() == "/" {
			m.showLogs = !m.showLogs
			return m, tea.ClearScreen
		}
		if x.String() == "ctrl+c" {
			m.closeAllPages()
			return m, tea.Quit
		}
	}

	curr := m.current()
	if curr.screen == nil {
		return m, tea.Quit
	}
	tr, cmd, err := curr.screen.Update(m.state, curr.page, msg)
	if err != nil {
		curr.page.SetError(err)
	}
	navCmd := m.applyTransition(tr)
	return m, tea.Batch(cmd, navCmd)
}

func (m *RootModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "Loading terminal UI..."
	}
	curr := m.current()
	if curr.screen == nil {
		return "No active screen"
	}

	rects := layout.Solve(m.width, m.height)
	help := curr.screen.Help(m.state, curr.page)
	help.Short = append(help.Short, "/: logs")
	body := curr.screen.View(m.state, curr.page, rects)
	main := lipgloss.NewStyle().
		Width(max(1, rects.Main.W)).
		Height(max(1, rects.Main.H)).
		MaxWidth(max(1, rects.Main.W)).
		MaxHeight(max(1, rects.Main.H)).
		Render(body)

	sidebarContent := m.sidebarContent(curr, rects)
	topHeader := m.styles.Header.Width(max(1, rects.Header.W)).Render(util.AppName)
	var content string
	if rects.Compact {
		header := m.styles.Header.Width(max(1, rects.Header.W)).Render(util.AppName + " | " + curr.page.Title)
		content = lipgloss.JoinVertical(lipgloss.Left, header, main)
	} else {
		sidebar := screens.Bounded(m.styles.Sidebar, rects.Sidebar.W, rects.Sidebar.H, sidebarContent)
		content = lipgloss.JoinVertical(lipgloss.Left, topHeader, lipgloss.JoinHorizontal(lipgloss.Top, main, sidebar))
	}

	editor := m.styles.Header.Width(max(1, rects.Editor.W)).Render(m.editorHint(curr))
	status := lipgloss.NewStyle().
		Width(max(1, rects.Status.W)).
		Height(max(1, rects.Status.H)).
		MaxWidth(max(1, rects.Status.W)).
		MaxHeight(max(1, rects.Status.H)).
		Render(components.RenderStatus(curr.page.Status, curr.page.Error, help.Short, m.state.ServerURL, m.state.ServerErr, rects.Status.W, m.styles))
	if m.showLogs {
		content = overlayBottom(content, m.logDrawer())
	}
	frameParts := []string{content}
	frameParts = append(frameParts, editor, status)
	frame := lipgloss.JoinVertical(lipgloss.Left, frameParts...)
	return m.styles.App.
		Width(max(1, m.width)).
		Height(max(1, m.height)).
		MaxWidth(max(1, m.width)).
		MaxHeight(max(1, m.height)).
		Render(frame)
}

func (m *RootModel) logDrawer() string {
	outerW := max(1, m.width)
	if outerW < 3 {
		return strings.Repeat(" ", outerW)
	}
	innerW := outerW - 2

	logLines := make([]string, 0, logDrawerLines)
	if m.state.LogTail != nil {
		for _, l := range m.state.LogTail(logDrawerLines) {
			logLines = append(logLines, singleLine(l))
		}
	}
	if len(logLines) == 0 {
		logLines = append(logLines, "(no logs yet)")
	}
	for len(logLines) < logDrawerLines {
		logLines = append(logLines, "")
	}
	if len(logLines) > logDrawerLines {
		logLines = logLines[len(logLines)-logDrawerLines:]
	}

	lines := make([]string, 0, logDrawerLines+2)
	lines = append(lines, logDrawerTopBorder(innerW))
	for _, l := range logLines {
		lines = append(lines, logDrawerBodyLine(l, innerW))
	}
	lines = append(lines, "╰"+strings.Repeat("─", innerW)+"╯")

	content := strings.Join(lines, "\n")
	drawerStyle := m.styles.Sidebar.
		Padding(0).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)
	return screens.Bounded(drawerStyle, outerW, len(lines), content)
}

func (m *RootModel) applyTransition(tr mvc.Transition) tea.Cmd {
	switch tr.Type {
	case mvc.TransitionStay:
		return nil
	case mvc.TransitionQuit:
		m.closeAllPages()
		return tea.Quit
	case mvc.TransitionPop:
		if len(m.stack) > 1 {
			m.current().page.Close()
			m.stack = m.stack[:len(m.stack)-1]
		}
		return nil
	case mvc.TransitionPush, mvc.TransitionRoute, mvc.TransitionReplace:
		scr, err := m.registry.Screen(tr.Screen)
		if err != nil {
			m.current().page.SetError(err)
			return nil
		}
		ps := mvc.NewPageState(m.state.Context, tr.Screen, strings.Title(tr.Screen), tr.Data, m.state.Logger)
		if tr.Type == mvc.TransitionReplace && len(m.stack) > 0 {
			m.current().page.Close()
			m.stack = m.stack[:len(m.stack)-1]
		}
		m.stack = append(m.stack, navEntry{screen: scr, page: ps})
		return scr.Init(m.state, ps)
	default:
		return nil
	}
}

func (m *RootModel) current() *navEntry {
	if len(m.stack) == 0 {
		return &navEntry{page: mvc.NewPageState(m.state.Context, "", "", nil, m.state.Logger)}
	}
	return &m.stack[len(m.stack)-1]
}

func (m *RootModel) closeAllPages() {
	for i := len(m.stack) - 1; i >= 0; i-- {
		m.stack[i].page.Close()
	}
}

func (m *RootModel) sidebarContent(curr *navEntry, rects layout.Rects) string {
	_ = rects
	lines := []string{fmt.Sprintf("section: %s", curr.page.Title)}
	if curr.page.Data != nil {
		for _, k := range curr.page.Data.Keys() {
			if k == "result" {
				continue
			}
			lines = append(lines, fmt.Sprintf("%s: %v", k, curr.page.Data[k]))
		}
	}
	if result := curr.page.EnsureData().GetStringArrayOpt("result"); len(result) > 0 {
		lines = append(lines, "", "recent output:")
		tail := result
		if len(tail) > 10 {
			tail = tail[len(tail)-10:]
		}
		lines = append(lines, tail...)
	}
	return strings.Join(lines, "\n")
}

func (m *RootModel) editorHint(curr *navEntry) string {
	return m.styles.Muted.Render(fmt.Sprintf("active: %s | use arrows + enter, b to go back", curr.page.Title))
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func logDrawerTopBorder(innerW int) string {
	if innerW < 1 {
		return "╮"
	}
	prefix := "─ Logs "
	if innerW <= lipgloss.Width(prefix) {
		return "╭" + truncateRunes(prefix, innerW) + "╮"
	}
	return "╭" + prefix + strings.Repeat("─", innerW-lipgloss.Width(prefix)) + "╮"
}

func logDrawerBodyLine(line string, innerW int) string {
	if innerW < 1 {
		return "││"
	}
	prefix := " "
	if innerW == 1 {
		return "│ │"
	}
	maxText := innerW - lipgloss.Width(prefix)
	text := truncateRunes(line, maxText)
	body := prefix + text
	pad := innerW - lipgloss.Width(body)
	if pad > 0 {
		body += strings.Repeat(" ", pad)
	}
	return "│" + body + "│"
}

func truncateRunes(s string, width int) string {
	if width < 1 {
		return ""
	}
	return ansi.Truncate(s, width, "…")
}

func overlayBottom(base string, overlay string) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")
	if len(baseLines) == 0 {
		return overlay
	}
	if len(overlayLines) >= len(baseLines) {
		return strings.Join(overlayLines[len(overlayLines)-len(baseLines):], "\n")
	}
	start := len(baseLines) - len(overlayLines)
	for i, l := range overlayLines {
		baseLines[start+i] = l
	}
	return strings.Join(baseLines, "\n")
}
