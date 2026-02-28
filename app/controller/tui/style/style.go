package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
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
	if th == nil || (th.Dark == nil && th.Light == nil) {
		th = theme.Default
	}
	c := colorsForMode(th)
	app := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Foreground)).Background(lipgloss.Color(c.Background))
	text := app
	sidebarBorder := "#111111"
	if IsDarkMode() {
		sidebarBorder = "#ffffff"
	}
	return Styles{
		App:    app,
		Header: text.Bold(true).Padding(0, 1).Background(lipgloss.Color(c.NavBackground)).Foreground(lipgloss.Color(c.NavForeground)),
		Panel: app.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(c.BackgroundMuted)).
			BorderBackground(lipgloss.Color(c.Background)).
			Padding(0, 1),
		Sidebar: app.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(sidebarBorder)).
			BorderBackground(lipgloss.Color(c.Background)).
			Padding(0, 1),
		Selected: text.Bold(true).Background(lipgloss.Color(c.NavBackground)).Foreground(lipgloss.Color(c.NavForeground)),
		Muted:    app.Foreground(lipgloss.Color(c.ForegroundMuted)),
		Status:   app.Foreground(lipgloss.Color(c.Success)).Padding(0, 1),
		Error:    app.Foreground(lipgloss.Color(c.Error)).Bold(true).Padding(0, 1),
	}
}

func colorsForMode(th *theme.Theme) *theme.Colors {
	dark := IsDarkMode()
	if dark && th.Dark != nil {
		return th.Dark
	}
	if !dark && th.Light != nil {
		return th.Light
	}
	if th.Dark != nil {
		return th.Dark
	}
	if th.Light != nil {
		return th.Light
	}
	return theme.Default.Dark
}

func IsDarkMode() bool {
	dark := lipgloss.HasDarkBackground()
	mode := strings.ToLower(strings.TrimSpace(util.GetEnv("tui_theme_mode", "auto")))
	switch mode {
	case "dark":
		return true
	case "light":
		return false
	default:
		return dark
	}
}
