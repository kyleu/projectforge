// $PF_HAS_MODULE(filesystem)$
package screens

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/util"
)

const (
	dataFileExplorerTitle = "file_explorer_title"
	dataFileExplorerRoot  = "file_explorer_root"
	dataFileExplorerStart = "file_explorer_start"
	dataFileExplorerPath  = "file_explorer_path"
	dataFileViewScroll    = "file_view_scroll"
	parentDir             = ".."
)

type FileBrowserScreen struct {
	picker filepicker.Model
	root   string
	title  string
}

type FileViewerScreen struct {
	path    string
	root    string
	title   string
	loading bool
	lines   []string
	meta    fileMeta
}

type fileViewLoadMsg struct {
	path  string
	meta  fileMeta
	lines []string
	err   error
}

type fileMeta struct {
	isDir      bool
	size       int64
	mode       os.FileMode
	modTime    time.Time
	childCount int
	dirCount   int
	fileCount  int
	lineCount  int
}

func NewFileBrowserScreen() *FileBrowserScreen {
	return &FileBrowserScreen{}
}

func (s *FileBrowserScreen) Key() string {
	return KeyFileBrowser
}

func (s *FileBrowserScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	root, title := fileExplorerContext(ps)
	start := fileExplorerStart(ps, root)
	s.root = root
	s.title = title
	ps.Title = title
	ps.SetStatus("Browse files")
	ps.Cursor = 0

	s.picker = filepicker.New()
	s.picker.CurrentDirectory = start
	s.picker.ShowHidden = true
	s.picker.ShowPermissions = true
	s.picker.ShowSize = true
	s.picker.DirAllowed = false
	s.picker.FileAllowed = true
	s.picker.AutoHeight = false
	// Keep "b" as screen-back; filepicker already handles h/left/backspace/esc for parent dir.
	s.picker.KeyMap.Back.SetKeys("h", "backspace", "left", "esc")

	return s.picker.Init()
}

func (s *FileBrowserScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if root, title := fileExplorerContext(ps); root != s.root || title != s.title {
		// Defensive reset if this singleton screen is reused with different page data.
		return mvc.Replace(KeyFileBrowser, ps.EnsureData().Clone()), nil, nil
	}

	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "b":
			return mvc.Pop(), nil, nil
		case "o":
			path := s.selectedPathForOpen(k)
			if path == "" {
				ps.SetStatus("No file or folder selected")
				return mvc.Stay(), nil, nil
			}
			path = clampToRoot(path, s.root)
			if err := OpenPath(path); err != nil {
				return mvc.Stay(), nil, errors.Wrapf(err, "unable to open path [%s]", path)
			}
			ps.SetStatus("Opened [%s]", filepath.Base(path))
			return mvc.Stay(), nil, nil
		case KeyEsc, KeyBackspace, KeyLeft, "h":
			if samePath(s.picker.CurrentDirectory, s.root) {
				return mvc.Pop(), nil, nil
			}
		}
	}

	var cmd tea.Cmd
	s.picker, cmd = s.picker.Update(msg)
	s.picker.CurrentDirectory = clampToRoot(s.picker.CurrentDirectory, s.root)
	if didSelect, path := s.picker.DidSelectFile(msg); didSelect {
		path = clampToRoot(path, s.root)
		return mvc.Push(KeyFileViewer, fileViewerData(s.title, s.root, path)), nil, nil
	}
	if didSelect, path := s.picker.DidSelectDisabledFile(msg); didSelect {
		ps.SetStatus("Cannot open [%s]", filepath.Base(path))
	}
	return mvc.Stay(), cmd, nil
}

func (s *FileBrowserScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	panelStyle := st.Panel
	contentW, contentH, _ := mainPanelContentSize(panelStyle, rects)
	bodyH := max(1, contentH-3)
	s.picker.SetHeight(bodyH)
	applyFilePickerStyles(&s.picker, st)

	meta := statPath(s.picker.CurrentDirectory)
	h1, h2 := browserHeaderLines(s.root, s.picker.CurrentDirectory, meta)
	body := s.renderBrowserPanelBody(st, contentW, bodyH, h1, h2)
	return renderScreenPanel(ps.Title, body, panelStyle, st, rects)
}

func (s *FileBrowserScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: move", "enter: open file/dir", "o: open in OS", "h/left: parent", "b: back"}}
}

func (s *FileBrowserScreen) selectedPathForOpen(k tea.KeyMsg) string {
	tmp := s.picker
	tmp.DirAllowed = true
	tmp.FileAllowed = true
	tmp.KeyMap.Open.SetKeys("o")
	tmp.KeyMap.Select.SetKeys("o")
	tmp, _ = tmp.Update(k)
	return tmp.Path
}

func (s *FileBrowserScreen) renderBrowserPanelBody(st style.Styles, width int, bodyH int, line1 string, line2 string) string {
	divider := st.Muted.Render(strings.Repeat("─", max(1, width-2)))
	listing := strings.TrimRight(s.picker.View(), "\n")
	if listing == "" {
		listing = " "
	}
	return strings.Join([]string{
		truncateLine(singleLine(line1), max(1, width)),
		truncateLine(singleLine(line2), max(1, width)),
		divider,
		fitVertical(listing, bodyH),
	}, "\n")
}

func NewFileViewerScreen() *FileViewerScreen {
	return &FileViewerScreen{}
}

func (s *FileViewerScreen) Key() string {
	return KeyFileViewer
}

func (s *FileViewerScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	root, title := fileExplorerContext(ps)
	path := clampToRoot(ps.EnsureData().GetStringOpt(dataFileExplorerPath), root)
	s.root = root
	s.title = title
	s.path = path
	s.loading = true
	s.lines = nil
	s.meta = fileMeta{}
	ps.Title = title
	ps.EnsureData()[dataFileViewScroll] = 0
	ps.SetStatus("Loading [%s]...", filepath.Base(path))
	return s.loadFileCmd(path)
}

func (s *FileViewerScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(fileViewLoadMsg); ok {
		if !samePath(m.path, s.path) {
			return mvc.Stay(), nil, nil
		}
		s.loading = false
		if m.err != nil {
			return mvc.Stay(), nil, m.err
		}
		s.lines = m.lines
		s.meta = m.meta
		ps.SetStatus("Loaded [%s]", filepath.Base(s.path))
		return mvc.Stay(), nil, nil
	}

	if delta, moved := menuMoveDelta(msg); moved {
		moveFileViewScroll(ps, delta)
		return mvc.Stay(), nil, nil
	}

	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "o":
			path := clampToRoot(s.path, s.root)
			if path == "" {
				ps.SetStatus("No file selected")
				return mvc.Stay(), nil, nil
			}
			if err := OpenPath(path); err != nil {
				return mvc.Stay(), nil, errors.Wrapf(err, "unable to open path [%s]", path)
			}
			ps.SetStatus("Opened [%s]", filepath.Base(path))
			return mvc.Stay(), nil, nil
		case KeyPgDown, "J":
			moveFileViewScroll(ps, 10)
			return mvc.Stay(), nil, nil
		case KeyPgUp, "K":
			moveFileViewScroll(ps, -10)
			return mvc.Stay(), nil, nil
		case KeyHome, "g":
			ps.EnsureData()[dataFileViewScroll] = 0
			return mvc.Stay(), nil, nil
		case KeyEnd, "G":
			if len(s.lines) > 0 {
				ps.EnsureData()[dataFileViewScroll] = len(s.lines) - 1
			}
			return mvc.Stay(), nil, nil
		case KeyEsc, KeyBackspace, KeyLeft, "h", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *FileViewerScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	newStyles := style.New(ts.Theme)
	panelStyle := newStyles.Panel
	contentW, contentH, _ := mainPanelContentSize(panelStyle, rects)
	bodyH := max(1, contentH-3)

	h1, h2 := viewerHeaderLines(s.root, s.path, s.meta)
	divider := newStyles.Muted.Render(strings.Repeat("─", max(1, contentW-2)))

	var content string
	if s.loading {
		content = "Loading..."
	} else {
		content = s.renderFileWindow(ps, bodyH)
	}

	body := strings.Join([]string{
		truncateLine(singleLine(h1), max(1, contentW)),
		truncateLine(singleLine(h2), max(1, contentW)),
		divider,
		fitVertical(content, bodyH),
	}, "\n")
	return renderScreenPanel(ps.Title, body, panelStyle, newStyles, rects)
}

func (s *FileViewerScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"up/down: scroll", "pgup/pgdn: page", "g/G: top/bottom", "o: open in OS", "b: back"}}
}

func (s *FileViewerScreen) renderFileWindow(ps *mvc.PageState, height int) string {
	if len(s.lines) == 0 {
		return "Empty file"
	}
	scroll := ps.EnsureData().GetIntOpt(dataFileViewScroll)
	if scroll < 0 {
		scroll = 0
	}
	maxOffset := max(0, len(s.lines)-max(1, height))
	if scroll > maxOffset {
		scroll = maxOffset
		ps.EnsureData()[dataFileViewScroll] = scroll
	}
	end := min(len(s.lines), scroll+max(1, height))
	return strings.Join(s.lines[scroll:end], "\n")
}

func (s *FileViewerScreen) loadFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		meta := statPath(path)
		if meta.isDir {
			return fileViewLoadMsg{path: path, meta: meta, err: errors.Errorf("can't view directory [%s]", path)}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fileViewLoadMsg{path: path, meta: meta, err: err}
		}
		lines, lineCount := renderHighlightedFile(path, data)
		meta.lineCount = lineCount
		return fileViewLoadMsg{path: path, meta: meta, lines: lines}
	}
}

func fileBrowserData(title string, root string, start string) util.ValueMap {
	root = cleanAbsPath(root)
	start = cleanAbsPath(start)
	if root == "" {
		root = "."
	}
	if start == "" {
		start = root
	}
	start = clampToRoot(start, root)
	return util.ValueMap{
		dataFileExplorerTitle: title,
		dataFileExplorerRoot:  root,
		dataFileExplorerStart: start,
	}
}

func fileViewerData(title string, root string, path string) util.ValueMap {
	ret := fileBrowserData(title, root, filepath.Dir(path))
	ret[dataFileExplorerPath] = cleanAbsPath(path)
	return ret
}

func fileExplorerContext(ps *mvc.PageState) (string, string) {
	d := ps.EnsureData()
	root := cleanAbsPath(d.GetStringOpt(dataFileExplorerRoot))
	if root == "" {
		root = "."
	}
	title := d.GetStringOpt(dataFileExplorerTitle)
	if title == "" {
		title = "Files"
	}
	return root, title
}

func fileExplorerStart(ps *mvc.PageState, root string) string {
	start := cleanAbsPath(ps.EnsureData().GetStringOpt(dataFileExplorerStart))
	if start == "" {
		start = root
	}
	return clampToRoot(start, root)
}

func cleanAbsPath(p string) string {
	if p == "" {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return filepath.Clean(p)
	}
	return filepath.Clean(abs)
}

func clampToRoot(path string, root string) string {
	path = cleanAbsPath(path)
	root = cleanAbsPath(root)
	if path == "" || root == "" {
		return path
	}
	if rel, err := filepath.Rel(root, path); err == nil && rel != parentDir && !strings.HasPrefix(rel, parentDir+string(os.PathSeparator)) {
		return path
	}
	return root
}

func samePath(a string, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func statPath(path string) fileMeta {
	info, err := os.Stat(path)
	if err != nil {
		return fileMeta{}
	}
	ret := fileMeta{isDir: info.IsDir(), size: info.Size(), mode: info.Mode(), modTime: info.ModTime()}
	if info.IsDir() {
		if entries, err := os.ReadDir(path); err == nil {
			ret.childCount = len(entries)
			for _, e := range entries {
				if e.IsDir() {
					ret.dirCount++
				} else {
					ret.fileCount++
				}
			}
		}
	}
	return ret
}

func browserHeaderLines(root string, curr string, meta fileMeta) (string, string) {
	rootLabel := displayPath(root, root)
	currLabel := displayPath(curr, root)
	line1 := fmt.Sprintf("Root: %s | Directory: %s", rootLabel, currLabel)
	line2 := fmt.Sprintf("Children: %d (%d dirs, %d files)", meta.childCount, meta.dirCount, meta.fileCount)
	if !meta.modTime.IsZero() {
		line2 += " | Modified: " + meta.modTime.Format(time.DateTime)
	}
	return line1, line2
}

func viewerHeaderLines(root string, path string, meta fileMeta) (string, string) {
	line1 := fmt.Sprintf("File: %s", displayPath(path, root))
	if meta.isDir {
		line2 := fmt.Sprintf("Directory | Children: %d (%d dirs, %d files)", meta.childCount, meta.dirCount, meta.fileCount)
		return line1, line2
	}
	size := util.ByteSizeSI(max(0, meta.size))
	line2 := fmt.Sprintf("Size: %s | Lines: %d | Mode: %s", size, meta.lineCount, meta.mode.String())
	if !meta.modTime.IsZero() {
		line2 += " | Modified: " + meta.modTime.Format(time.DateTime)
	}
	return line1, line2
}

func displayPath(path string, root string) string {
	if path == "" {
		return ""
	}
	if root == "" {
		return path
	}
	rel, err := filepath.Rel(root, path)
	if err == nil && rel == "." {
		return "."
	}
	if err == nil && rel != parentDir && !strings.HasPrefix(rel, parentDir+string(os.PathSeparator)) {
		return filepath.Join(".", rel)
	}
	return path
}

func moveFileViewScroll(ps *mvc.PageState, delta int) {
	if delta == 0 {
		return
	}
	d := ps.EnsureData()
	d[dataFileViewScroll] = d.GetIntOpt(dataFileViewScroll) + delta
	if d.GetIntOpt(dataFileViewScroll) < 0 {
		d[dataFileViewScroll] = 0
	}
}

func fitVertical(s string, height int) string {
	if height <= 1 {
		if s == "" {
			return ""
		}
		if i := strings.IndexByte(s, '\n'); i >= 0 {
			return s[:i]
		}
		return s
	}
	lines := strings.Split(s, "\n")
	if len(lines) >= height {
		return strings.Join(lines[:height], "\n")
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func applyFilePickerStyles(p *filepicker.Model, st style.Styles) {
	p.Styles.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	p.Styles.Selected = lipgloss.NewStyle().Bold(true)
	if style.IsDarkMode() {
		p.Styles.Directory = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
		p.Styles.File = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		p.Styles.Permission = st.Muted
		p.Styles.FileSize = st.Muted.Width(7).Align(lipgloss.Right)
		return
	}
	p.Styles.Directory = lipgloss.NewStyle().Foreground(lipgloss.Color("27"))
	p.Styles.File = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	p.Styles.Permission = st.Muted
	p.Styles.FileSize = st.Muted.Width(7).Align(lipgloss.Right)
}

func renderHighlightedFile(path string, data []byte) ([]string, int) {
	if len(data) == 0 {
		return []string{" "}, 0
	}
	if !utf8.Valid(data) {
		msg := fmt.Sprintf("Binary file (%s) cannot be syntax-highlighted", util.ByteSizeSI(int64(len(data))))
		return []string{msg}, 0
	}
	content := string(data)
	lineCount := strings.Count(content, "\n")
	if !strings.HasSuffix(content, "\n") {
		lineCount++
	}
	if lineCount < 0 {
		lineCount = 0
	}

	lex := lexers.Match(path)
	if lex == nil {
		lex = lexers.Analyse(content) //nolint:misspell // upstream API uses British spelling
	}
	if lex == nil {
		lex = lexers.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	styleName := "github"
	if style.IsDarkMode() {
		styleName = "monokai"
	}
	chromaStyle := styles.Get(styleName)
	if chromaStyle == nil {
		chromaStyle = styles.Fallback
	}

	var buf bytes.Buffer
	it, err := lex.Tokenise(nil, content)
	if err != nil {
		return strings.Split(content, "\n"), lineCount
	}
	if err = formatter.Format(&buf, chromaStyle, it); err != nil {
		return strings.Split(content, "\n"), lineCount
	}
	out := strings.TrimRight(buf.String(), "\n")
	if out == "" {
		return []string{" "}, lineCount
	}
	return strings.Split(out, "\n"), lineCount
}
