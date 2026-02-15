package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
)

func RenderMenuList(items menu.Items, cursor int, st style.Styles, width int) string {
	if len(items) == 0 {
		return st.Muted.Render("No items available")
	}
	if width < 8 {
		width = 8
	}
	titleStyle := st.App.Bold(true)
	descStyle := st.Muted

	out := make([]string, 0, len(items))
	for i, item := range items {
		line := renderMenuRow(item.Title, item.Description, width, titleStyle, descStyle)
		if i == cursor {
			line = st.Selected.Render(stripANSI(line))
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func renderMenuRow(title string, desc string, width int, titleStyle lipgloss.Style, descStyle lipgloss.Style) string {
	const (
		prefix = "  "
		sep    = " - "
	)
	available := width - runeLen(prefix)
	if available < 1 {
		return prefix
	}

	title = singleLine(title)
	desc = singleLine(desc)
	titleOut := truncateEllipsis(title, available)
	descOut := ""
	used := runeLen(titleOut)
	remaining := available - used
	if desc != "" && remaining > runeLen(sep)+2 {
		descOut = truncateEllipsis(desc, remaining-runeLen(sep))
	}

	ret := prefix + titleStyle.Render(titleOut)
	if descOut != "" {
		ret += sep + descStyle.Render(descOut)
	}
	return ret
}

func truncateEllipsis(s string, width int) string {
	if width < 1 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width <= 1 {
		return string(r[:width])
	}
	return string(r[:width-1]) + "…"
}

func singleLine(s string) string {
	// Collapse all whitespace (including newlines/tabs) so each row is always one visual line.
	return strings.Join(strings.Fields(s), " ")
}

func runeLen(s string) int {
	return len([]rune(s))
}

func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inEscape {
			if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				inEscape = false
			}
			continue
		}
		if ch == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}
