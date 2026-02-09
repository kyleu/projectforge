package tui

const (
	tuiKeyAll       = "all"
	tuiKeyBack      = "back"
	tuiKeyEsc       = "esc"
	tuiKeyUp        = "up"
	tuiKeyDown      = "down"
	tuiKeyEnter     = "enter"
	tuiKeyCtrlC     = "ctrl+c"
	tuiCursorSpacer = "  "
)

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
