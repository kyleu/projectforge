package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/doctor"
)

func handleMessage(t *TUI, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case projectsLoadedMsg:
		t.Config.projectsLoading = false
		t.Config.projectsErr = msg.err
		prjs := projectsFor(t)
		if t.Screen == screenProjects && t.Screen.Cursor() >= len(prjs) {
			t.Screen.ResetCursor()
		}
	case doctorChecksLoadedMsg:
		t.Config.doctorLoading = false
		t.Config.doctorErr = msg.err
		t.Config.doctorChecks = msg.checks
		t.Config.doctorResults = map[string]*doctor.Result{}
		if t.Screen == screenDoctor && t.Screen.Cursor() >= len(t.Config.doctorChecks) {
			t.Screen.ResetCursor()
		}
	case doctorCheckResultMsg:
		t.Config.doctorRunning = false
		t.Config.doctorErr = msg.err
		if msg.result != nil {
			t.Config.doctorResults[msg.result.Key] = msg.result
		}
	case doctorAllResultsMsg:
		t.Config.doctorRunning = false
		t.Config.doctorErr = msg.err
		for _, r := range msg.results {
			if r != nil {
				t.Config.doctorResults[r.Key] = r
			}
		}
	}
	return t, nil
}
