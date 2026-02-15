package style

import (
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/lib/theme"
)

type Styles struct {
	App      lipgloss.Style
	Header   lipgloss.Style
	Panel    lipgloss.Style
	Sidebar  lipgloss.Style
	Selected lipgloss.Style
	Muted    lipgloss.Style
	Status   lipgloss.Style
	Error    lipgloss.Style
}

func New(th *theme.Theme) Styles {
	if th == nil || th.Dark == nil {
		th = theme.Default
	}
	c := th.Dark
	text := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Foreground))
	app := text.Background(lipgloss.Color(c.Background))
	return Styles{
		App:      app,
		Header:   text.Bold(true).Padding(0, 1).Background(lipgloss.Color(c.NavBackground)).Foreground(lipgloss.Color(c.NavForeground)),
		Panel:    text.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(c.BackgroundMuted)).Padding(0, 1),
		Sidebar:  text.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(c.Border)).Padding(0, 1),
		Selected: text.Bold(true).Background(lipgloss.Color(c.NavBackground)).Foreground(lipgloss.Color(c.NavForeground)),
		Muted:    text.Foreground(lipgloss.Color(c.ForegroundMuted)),
		Status:   text.Foreground(lipgloss.Color(c.Success)).Padding(0, 1),
		Error:    text.Foreground(lipgloss.Color(c.Error)).Bold(true).Padding(0, 1),
	}
}
