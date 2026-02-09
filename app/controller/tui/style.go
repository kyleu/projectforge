package tui

import "github.com/charmbracelet/lipgloss"

const (
	colorPrimary    = "#0a2638"
	colorAccent     = "#50677d"
	colorHighlight  = "#cdd9e3"
	colorTextLight  = "#eeeeee"
	colorTextMuted  = "#93a1af"
	colorBackground = "#121212"
)

var (
	hl   = lipgloss.Color(colorHighlight)
	acct = lipgloss.Color(colorAccent)
	tl   = lipgloss.Color(colorTextLight)
	tm   = lipgloss.Color(colorTextMuted)

	titleStyle          = lipgloss.NewStyle().Foreground(tl).Background(lipgloss.Color(colorPrimary)).Bold(true).Padding(1, 4)
	splashTitleStyle    = lipgloss.NewStyle().Foreground(hl).Bold(true)
	splashSubtitleStyle = lipgloss.NewStyle().Foreground(acct).Italic(true)
	pressKeyStyle       = lipgloss.NewStyle().Foreground(tm).Blink(true)
	menuItemStyle       = lipgloss.NewStyle().Foreground(tl).PaddingLeft(2)
	menuSelectedStyle   = lipgloss.NewStyle().Foreground(tl).Background(acct).Bold(true).PaddingLeft(2)
	menuCursorStyle     = lipgloss.NewStyle().Foreground(hl).Bold(true)
	resultStyle         = lipgloss.NewStyle().Foreground(hl).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(acct).Padding(1, 3)
	helpStyle           = lipgloss.NewStyle().Foreground(tm).MarginTop(1)
	containerStyle      = lipgloss.NewStyle().Align(lipgloss.Center).AlignVertical(lipgloss.Center)
	logDividerStyle     = lipgloss.NewStyle().Foreground(acct)
	logHeaderStyle      = lipgloss.NewStyle().Foreground(tm).Bold(true)
	logLineStyle        = lipgloss.NewStyle().Foreground(tl)
	statusLineStyle     = lipgloss.NewStyle().Foreground(tl).Background(lipgloss.Color(colorPrimary)).Padding(0, 1)
)
