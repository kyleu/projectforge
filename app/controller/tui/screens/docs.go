package screens

import (
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/cmenu"
	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/doc"
)

type DocumentationScreen struct {
	root       *menu.Item
	stack      []*menu.Item
	cursorPath []int
	activeFile string
	scroll     int
	loading    bool
	viewportW  int
	lines      []string
}

type docRenderMsg struct {
	path  string
	lines []string
	err   error
}

func NewDocumentationScreen() *DocumentationScreen {
	return &DocumentationScreen{}
}

func (s *DocumentationScreen) Key() string {
	return KeyDocs
}

func (s *DocumentationScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	s.root = cmenu.DocsMenu(ts.Logger)
	s.stack = nil
	s.cursorPath = []int{0}
	s.activeFile = ""
	s.scroll = 0
	s.loading = false
	s.viewportW = 0
	s.lines = nil
	ps.Cursor = 0
	syncDocsTitle(ps, nil)
	ps.SetStatus("Browse embedded documentation")
	return nil
}

func (s *DocumentationScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	items := s.currentItems()
	ps.Cursor = clampMenuCursor(ps.Cursor, len(items))

	if m, ok := msg.(docRenderMsg); ok {
		return s.updateRenderMsg(ps, m)
	}

	if s.activeFile == "" {
		if delta, moved := menuMoveDelta(msg); moved {
			ps.Cursor = moveMenuCursor(ps.Cursor, len(items), delta)
			s.cursorPath[len(s.cursorPath)-1] = ps.Cursor
			return mvc.Stay(), nil, nil
		}
	} else if delta, moved := menuMoveDelta(msg); moved {
		s.moveScroll(delta)
		return mvc.Stay(), nil, nil
	}

	if m, ok := msg.(tea.KeyMsg); ok {
		return s.updateKeyMsg(ps, items, m)
	}
	return mvc.Stay(), nil, nil
}

func (s *DocumentationScreen) updateRenderMsg(ps *mvc.PageState, m docRenderMsg) (mvc.Transition, tea.Cmd, error) {
	if m.path != s.activeFile {
		return mvc.Stay(), nil, nil
	}
	s.loading = false
	if m.err != nil {
		return mvc.Stay(), nil, m.err
	}
	s.lines = m.lines
	ps.SetStatus("File loaded")
	return mvc.Stay(), nil, nil
}

func (s *DocumentationScreen) updateKeyMsg(ps *mvc.PageState, items menu.Items, m tea.KeyMsg) (mvc.Transition, tea.Cmd, error) {
	switch m.String() {
	case KeyEnter:
		return s.openSelected(ps, items)
	case KeyEsc, KeyBackspace, "b", KeyLeft, "h":
		return s.goBack(ps)
	default:
		return mvc.Stay(), nil, nil
	}
}

func (s *DocumentationScreen) openSelected(ps *mvc.PageState, items menu.Items) (mvc.Transition, tea.Cmd, error) {
	if s.activeFile != "" || len(items) == 0 {
		return mvc.Stay(), nil, nil
	}
	selected := items[ps.Cursor]
	if len(selected.Children) > 0 || selected.Route == "" {
		s.stack = append(s.stack, selected)
		s.cursorPath = append(s.cursorPath, 0)
		ps.Cursor = 0
		syncDocsTitle(ps, s.stack)
		ps.SetStatus("Folder: %s", folderLabel(s.stack))
		return mvc.Stay(), nil, nil
	}
	docPath, err := docPathFromRoute(selected.Route)
	if err != nil {
		return mvc.Stay(), nil, err
	}
	s.activeFile = docPath
	s.scroll = 0
	s.lines = nil
	syncDocsTitle(ps, s.stack)
	ps.SetStatus("File: %s (loading...)", selected.Title)
	s.loading = true
	return mvc.Stay(), s.renderMarkdownCmd(docPath), nil
}

func (s *DocumentationScreen) goBack(ps *mvc.PageState) (mvc.Transition, tea.Cmd, error) {
	if s.activeFile != "" {
		s.activeFile = ""
		s.scroll = 0
		s.loading = false
		s.lines = nil
		syncDocsTitle(ps, s.stack)
		ps.SetStatus("Folder: %s", folderLabel(s.stack))
		return mvc.Stay(), nil, nil
	}
	if len(s.stack) > 0 {
		s.stack = s.stack[:len(s.stack)-1]
		s.cursorPath = s.cursorPath[:len(s.cursorPath)-1]
		ps.Cursor = s.cursorPath[len(s.cursorPath)-1]
		syncDocsTitle(ps, s.stack)
		ps.SetStatus("Folder: %s", folderLabel(s.stack))
		return mvc.Stay(), nil, nil
	}
	return mvc.Pop(), nil, nil
}

func (s *DocumentationScreen) SidebarContent(ts *mvc.State, ps *mvc.PageState, _ layout.Rects) (string, bool) {
	styles := style.New(ts.Theme)
	lines := []string{"Documentation", ""}
	if s.activeFile == "" {
		items := s.currentItems()
		cursor := clampMenuCursor(ps.Cursor, len(items))
		lines = AppendSidebarProp(lines, styles, "folder", folderLabel(s.stack))
		lines = AppendSidebarProp(lines, styles, "items", len(items))
		if len(items) > 0 {
			item := items[cursor]
			kind := "file"
			if len(item.Children) > 0 || item.Route == "" {
				kind = "folder"
			}
			lines = AppendSidebarProp(lines, styles, "selected", item.Title)
			lines = AppendSidebarProp(lines, styles, "type", kind)
		}
		lines = append(lines, "", "keys:", "enter open", "b back")
	} else {
		lines = AppendSidebarProp(lines, styles, "file", s.activeFile)
		if s.loading {
			lines = AppendSidebarProp(lines, styles, "status", "loading...")
		} else {
			lines = AppendSidebarProp(lines, styles, "lines", len(s.lines))
			if len(s.lines) > 0 {
				lines = AppendSidebarProp(lines, styles, "at", min(len(s.lines), s.scroll+1))
			}
		}
		lines = append(lines, "", "keys:", "up/down scroll", "b back")
	}
	return strings.Join(lines, "\n"), true
}

func (s *DocumentationScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	items := s.currentItems()
	ps.Cursor = clampMenuCursor(ps.Cursor, len(items))
	panelStyle := styles.Panel
	contentW, contentH, _ := mainPanelContentSize(panelStyle, rects)
	s.viewportW = contentW

	var body string
	if s.activeFile == "" {
		body = components.RenderMenuList(items, ps.Cursor, styles, contentW)
	} else {
		body = s.renderMarkdownWindow(contentH)
	}

	return renderScreenPanel(ps.Title, body, panelStyle, styles, rects)
}

func (s *DocumentationScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	if s.activeFile != "" {
		return HelpBindings{Short: []string{"up/down: scroll", "left|b: back"}}
	}
	return HelpBindings{Short: []string{"up/down: move", "enter: open", "left|b: back"}}
}

func (s *DocumentationScreen) currentItems() menu.Items {
	curr := s.root
	for _, n := range s.stack {
		curr = n
	}
	if curr == nil {
		return nil
	}
	return curr.Children.Visible()
}

func (s *DocumentationScreen) renderMarkdownWindow(height int) string {
	if s.loading {
		return ""
	}
	if s.activeFile == "" {
		return ""
	}
	if len(s.lines) == 0 {
		return "No content"
	}
	maxOffset := max(0, len(s.lines)-max(1, height))
	if s.scroll < 0 {
		s.scroll = 0
	}
	if s.scroll > maxOffset {
		s.scroll = maxOffset
	}
	end := min(len(s.lines), s.scroll+max(1, height))
	return strings.Join(s.lines[s.scroll:end], "\n")
}

func (s *DocumentationScreen) moveScroll(delta int) {
	if delta == 0 {
		return
	}
	s.scroll += delta
	if s.scroll < 0 {
		s.scroll = 0
	}
}

func (s *DocumentationScreen) renderMarkdownCmd(docPath string) tea.Cmd {
	return func() tea.Msg {
		data, err := doc.Content(docPath)
		if err != nil {
			return docRenderMsg{path: docPath, err: err}
		}
		lines, _ := renderHighlightedFile(docPath, data)
		return docRenderMsg{path: docPath, lines: lines}
	}
}

func syncDocsTitle(ps *mvc.PageState, stack []*menu.Item) {
	if len(stack) == 0 {
		ps.Title = "Documentation"
		return
	}
	parts := make([]string, 0, len(stack)+1)
	parts = append(parts, "Documentation")
	for _, n := range stack {
		parts = append(parts, n.Title)
	}
	ps.Title = strings.Join(parts, " / ")
}

func folderLabel(stack []*menu.Item) string {
	if len(stack) == 0 {
		return "/"
	}
	parts := make([]string, 0, len(stack))
	for _, n := range stack {
		parts = append(parts, n.Title)
	}
	return path.Join(parts...)
}

func docPathFromRoute(route string) (string, error) {
	const prefix = "/docs/"
	if !strings.HasPrefix(route, prefix) {
		return "", errors.Errorf("route [%s] is not a docs route", route)
	}
	pth := strings.TrimPrefix(route, prefix)
	if pth == "" {
		return "", errors.New("cannot use empty path")
	}
	return pth + ".md", nil
}
