package screens

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/util"
)

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
		return util.KeyEllipsis
	}
	return string(r[:width-1]) + util.KeyEllipsis
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func runeLen(s string) int {
	return len([]rune(s))
}

func mainPanelOuterHeight(rects layout.Rects) int {
	if rects.Compact {
		return max(1, rects.Main.H-1)
	}
	return max(1, rects.Main.H)
}

func mainPanelContentSize(panelStyle lipgloss.Style, rects layout.Rects) (int, int, int) {
	outerH := mainPanelOuterHeight(rects)
	contentW, contentH := panelContentSize(panelStyle, rects.Main.W, outerH)
	return contentW, contentH, outerH
}

func renderScreenPanel(title string, body string, panelStyle lipgloss.Style, styles style.Styles, rects layout.Rects) string {
	_, _, outerH := mainPanelContentSize(panelStyle, rects)
	panel := renderTextPanel(body, panelStyle, rects.Main.W, outerH)
	return renderScreenFrame(title, panel, styles, rects)
}

func renderScreenFrame(title string, content string, styles style.Styles, rects layout.Rects) string {
	if !rects.Compact {
		return content
	}
	header := styles.Header.Width(max(1, rects.Main.W)).Render(title)
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func ContentSize(st lipgloss.Style, outerW int, outerH int) (int, int) {
	w := max(1, outerW-st.GetHorizontalFrameSize())
	h := max(1, outerH-st.GetVerticalFrameSize())
	return w, h
}

func Bounded(st lipgloss.Style, outerW int, outerH int, content string) string {
	ow := max(1, outerW)
	oh := max(1, outerH)
	cw, ch := ContentSize(st, ow, oh)
	return st.
		Width(cw).
		Height(ch).
		MaxWidth(ow).
		MaxHeight(oh).
		Render(content)
}

func AppendSidebarProp(lines []string, styles style.Styles, key string, value any) []string {
	label := styles.Muted.Bold(true).Render(key)
	return append(lines, label, fmt.Sprint(value), "")
}

func OpenInBrowser(url string) error {
	var cmd *exec.Cmd
	ctx := context.Background()
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "open", url)
	case "windows":
		cmd = exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.CommandContext(ctx, "xdg-open", url)
	}
	return cmd.Start()
}

func OpenPath(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
