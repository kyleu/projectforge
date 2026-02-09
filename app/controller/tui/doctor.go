package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

var screenDoctor = NewScreen(
	"doctor", "Doctor", "checks system health", renderDoctor,
	`"esc": back`, `"↑"/"↓": move`, `"enter": run check`, `"a": run all`, `"r": reload`, `"q": quit`,
)

type doctorChecksLoadedMsg struct {
	checks doctor.Checks
	err    error
}

type doctorCheckResultMsg struct {
	result *doctor.Result
	err    error
}

type doctorAllResultsMsg struct {
	results doctor.Results
	err     error
}

func renderDoctor(t *TUI) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Doctor"))
	b.WriteString("\n\n")

	switch {
	case t.Config.doctorLoading:
		b.WriteString("Loading checks...")
	case t.Config.doctorErr != nil:
		b.WriteString("Doctor error:\n")
		b.WriteString(t.Config.doctorErr.Error())
	case len(t.Config.doctorChecks) == 0:
		b.WriteString("No doctor checks available for the current projects.")
	default:
		b.WriteString(renderDoctorChecks(t))
		b.WriteString("\n")
		if t.Config.doctorRunning {
			b.WriteString(helpStyle.Render("Running check..."))
			b.WriteString("\n\n")
		}
		cursor := t.Screen.Cursor()
		if cursor >= 0 && cursor < len(t.Config.doctorChecks) {
			b.WriteString(resultStyle.Render(doctorDetails(t.Config.doctorChecks[cursor], t.Config.doctorResults[t.Config.doctorChecks[cursor].Key])))
		}
	}

	return containerStyle.Width(t.width).Height(t.height).Render(b.String())
}

func renderDoctorChecks(t *TUI) string {
	items := make(menu.Items, 0, len(t.Config.doctorChecks))
	for _, c := range t.Config.doctorChecks {
		label := c.Title
		if r, ok := t.Config.doctorResults[c.Key]; ok {
			status := "ok"
			if r.Status == util.KeyError {
				status = "error"
			}
			label = fmt.Sprintf("%s [%s]", c.Title, status)
		}
		items = append(items, &menu.Item{Key: c.Key, Title: label})
	}
	return RenderMenuOptions(t.Screen.Cursor(), items)
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
	if t.Config.doctorLoading || t.Config.doctorRunning {
		return nil
	}

	switch key {
	case tuiKeyUp, "k":
		if t.Screen.Cursor() > 0 {
			t.Screen.ModifyCursor(-1)
		}
	case tuiKeyDown, "j":
		if t.Screen.Cursor() < len(t.Config.doctorChecks)-1 {
			t.Screen.ModifyCursor(1)
		}
	case "r":
		t.Config.doctorLoading = true
		t.Config.doctorErr = nil
		t.Config.doctorResults = map[string]*doctor.Result{}
		return loadDoctorChecksCmd(t)
	case "a":
		t.Config.doctorRunning = true
		t.Config.doctorErr = nil
		return runDoctorAllCmd(t)
	case tuiKeyEnter, " ":
		cursor := t.Screen.Cursor()
		if cursor < 0 || cursor >= len(t.Config.doctorChecks) {
			return nil
		}
		t.Config.doctorRunning = true
		t.Config.doctorErr = nil
		return runDoctorCheckCmd(t, t.Config.doctorChecks[cursor].Key)
	}
	return nil
}

func loadDoctorChecksCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		prjs, err := t.st.Services.Projects.Refresh(t.logger)
		if err != nil {
			return doctorChecksLoadedMsg{err: errors.Wrap(err, "unable to load projects")}
		}
		checks.SetModules(t.st.Services.Modules)
		ret := checks.ForModules(prjs.AllModules())
		return doctorChecksLoadedMsg{checks: ret}
	}
}

func runDoctorCheckCmd(t *TUI, key string) tea.Cmd {
	return func() tea.Msg {
		c := checks.GetCheck(key)
		if c == nil {
			return doctorCheckResultMsg{err: errors.Errorf("no doctor check found for key [%s]", key)}
		}
		res := c.Check(t.ctx, t.logger)
		if res == nil {
			return doctorCheckResultMsg{err: errors.Errorf("doctor check [%s] does not apply to this platform", key)}
		}
		return doctorCheckResultMsg{result: res}
	}
}

func runDoctorAllCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		modules := projectsFor(t).AllModules()
		if len(modules) == 0 {
			prjs, err := t.st.Services.Projects.Refresh(t.logger)
			if err != nil {
				return doctorAllResultsMsg{err: errors.Wrap(err, "unable to load projects")}
			}
			modules = prjs.AllModules()
		}
		checks.SetModules(t.st.Services.Modules)
		ret := checks.CheckAll(t.ctx, modules, t.logger)
		return doctorAllResultsMsg{results: ret}
	}
}
