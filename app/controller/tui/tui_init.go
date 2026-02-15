package tui

func InitTUI(t *TUI) (*TUI, error) {
	_, _ = t.st.Services.Projects.Refresh(t.logger)
	return t, nil
}
