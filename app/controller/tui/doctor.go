package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

var screenDoctor = &Screen{
	Key:     "doctor",
	Title:   "Doctor",
	Hotkeys: []string{`"esc": back`, `"↑"/"↓": move`, `"enter": run check`, `"a": run all`, `"r": reload`, `"q": quit`},
	Render:  renderDoctor,
}

func renderDoctor(t *TUI) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Doctor"))
	b.WriteString("\n\n")

	switch {
	case t.doctorLoading:
		b.WriteString("Loading checks...")
	case t.doctorErr != nil:
		b.WriteString("Doctor error:\n")
		b.WriteString(t.doctorErr.Error())
	case len(t.doctorChecks) == 0:
		b.WriteString("No doctor checks available for the current projects.")
	default:
		b.WriteString(renderDoctorChecks(t))
		b.WriteString("\n")
		if t.doctorRunning {
			b.WriteString(helpStyle.Render("Running check..."))
			b.WriteString("\n\n")
		}
		if t.Screen.Cursor >= 0 && t.Screen.Cursor < len(t.doctorChecks) {
			b.WriteString(resultStyle.Render(doctorDetails(t.doctorChecks[t.Screen.Cursor], t.doctorResults[t.doctorChecks[t.Screen.Cursor].Key])))
		}
	}

	return containerStyle.Width(t.width).Height(t.height).Render(b.String())
}

func renderDoctorChecks(t *TUI) string {
	items := make(menu.Items, 0, len(t.doctorChecks))
	for _, c := range t.doctorChecks {
		label := c.Title
		if r, ok := t.doctorResults[c.Key]; ok {
			status := "ok"
			if r.Status == util.KeyError {
				status = "error"
			}
			label = fmt.Sprintf("%s [%s]", c.Title, status)
		}
		items = append(items, &menu.Item{Key: c.Key, Title: label})
	}
	return RenderMenuOptions(t.Screen.Cursor, items)
}

func doctorDetails(c *doctor.Check, r *doctor.Result) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Check: %s\n", c.Title))
	b.WriteString(fmt.Sprintf("Summary: %s\n", orDash(c.Summary)))
	b.WriteString(fmt.Sprintf("Section: %s\n", orDash(c.Section)))
	if len(c.Modules) > 0 {
		b.WriteString(fmt.Sprintf("Modules: %s\n", strings.Join(c.Modules, ", ")))
	}
	if r == nil {
		b.WriteString("\nNo result yet. Press Enter to run this check.")
		return b.String()
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Result: %s\n", strings.ToUpper(r.Status)))
	b.WriteString(fmt.Sprintf("Duration: %s\n", util.MicrosToMillis(r.Duration)))
	if len(r.Errors) > 0 {
		b.WriteString("Errors:\n")
		for _, e := range r.Errors {
			b.WriteString("- ")
			b.WriteString(e.String())
			b.WriteByte('\n')
		}
	}
	if len(r.CleanSolutions()) > 0 {
		b.WriteString("Solutions:\n")
		for _, s := range r.CleanSolutions() {
			b.WriteString("- ")
			b.WriteString(s)
			b.WriteByte('\n')
		}
	}
	return strings.TrimSpace(b.String())
}

func onKeyDoctor(key string, t *TUI) tea.Cmd {
	if key == "q" {
		t.quitting = true
		return tea.Quit
	}
	if key == tuiKeyEsc {
		t.Screen = screenMenu
		return nil
	}
	if t.doctorLoading || t.doctorRunning {
		return nil
	}

	switch key {
	case tuiKeyUp, "k":
		if t.Screen.Cursor > 0 {
			t.Screen.Cursor--
		}
	case tuiKeyDown, "j":
		if t.Screen.Cursor < len(t.doctorChecks)-1 {
			t.Screen.Cursor++
		}
	case "r":
		t.doctorLoading = true
		t.doctorErr = nil
		t.doctorResults = map[string]*doctor.Result{}
		return loadDoctorChecksCmd(t)
	case "a":
		t.doctorRunning = true
		t.doctorErr = nil
		return runDoctorAllCmd(t)
	case tuiKeyEnter, " ":
		if t.Screen.Cursor < 0 || t.Screen.Cursor >= len(t.doctorChecks) {
			return nil
		}
		t.doctorRunning = true
		t.doctorErr = nil
		return runDoctorCheckCmd(t, t.doctorChecks[t.Screen.Cursor].Key)
	}
	return nil
}
