package tui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/util"
)

func (t *TUI) AddLog(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) {
	line := formatLogLine(level, occurred, loggerName, message, caller, stack, fields)
	t.logsMu.Lock()
	defer t.logsMu.Unlock()
	t.logs = append(t.logs, line)
	if len(t.logs) > 200 {
		t.logs = t.logs[1:]
	}
}

func (t *TUI) LastLogs(limit int) []string {
	t.logsMu.RLock()
	defer t.logsMu.RUnlock()
	if limit <= 0 || len(t.logs) == 0 {
		return nil
	}
	if limit > len(t.logs) {
		limit = len(t.logs)
	}
	start := len(t.logs) - limit
	ret := make([]string, 0, limit)
	for _, line := range t.logs[start:] {
		ret = append(ret, line)
	}
	return ret
}

var (
	logTimeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	logLoggerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	logCallerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	logFieldsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	logStackStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
)

func formatLogLine(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) string {
	levelText := strings.ToUpper(strings.TrimSpace(level))
	if levelText == "" {
		levelText = "INFO"
	}
	pieces := []string{
		logTimeStyle.Render(occurred.Format("15:04:05.000")),
		renderLevel(levelText),
	}

	loggerName = strings.TrimSpace(loggerName)
	if loggerName != "" {
		pieces = append(pieces, logLoggerStyle.Render("["+loggerName+"]"))
	}

	msg := inlineText(message)
	if msg != "" {
		pieces = append(pieces, msg)
	}

	if c := shortCaller(caller); c != "" {
		pieces = append(pieces, logCallerStyle.Render(c))
	}
	if f := summarizeFields(fields); f != "" {
		pieces = append(pieces, logFieldsStyle.Render(f))
	}
	if s := summarizeStack(stack); s != "" {
		pieces = append(pieces, logStackStyle.Render(s))
	}

	return strings.Join(pieces, " ")
}

func renderLevel(level string) string {
	base := lipgloss.NewStyle().Bold(true)
	switch strings.ToLower(level) {
	case "debug":
		base = base.Foreground(lipgloss.Color("141"))
	case "warn":
		base = base.Foreground(lipgloss.Color("214"))
	case "error", "dpanic", "panic", "fatal":
		base = base.Foreground(lipgloss.Color("196"))
	default:
		base = base.Foreground(lipgloss.Color("42"))
	}
	return base.Render(fmt.Sprintf("%-5s", level))
}

func shortCaller(caller util.ValueMap) string {
	if len(caller) == 0 {
		return ""
	}
	file := strings.TrimSpace(caller.GetStringOpt("file"))
	line := caller.GetIntOpt("line")
	if line == 0 {
		if x, err := strconv.Atoi(strings.TrimSpace(caller.GetStringOpt("line"))); err == nil {
			line = x
		}
	}
	if file == "" {
		return ""
	}
	file = filepath.Base(file)
	if line > 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return file
}

func summarizeFields(fields util.ValueMap) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for _, k := range fields.Keys() {
		v := inlineText(fmt.Sprint(fields[k]))
		parts = append(parts, fmt.Sprintf("%s=%s", k, truncateText(v, 24)))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

func summarizeStack(stack string) string {
	lines := util.StringSplitLines(stack)
	for _, line := range lines {
		line = inlineText(line)
		if line == "" {
			continue
		}
		return "stack: " + truncateText(line, 48)
	}
	return ""
}

func inlineText(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func truncateText(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	if limit == 1 {
		return "…"
	}
	return string(r[:limit-1]) + "…"
}
