package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/project"
)

type projectsLoadedMsg struct {
	projects project.Projects
	err      error
}

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

func loadProjectsCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		prjs, err := t.st.Services.Projects.Refresh(t.logger)
		return projectsLoadedMsg{projects: prjs, err: err}
	}
}

func loadDoctorChecksCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		prjs, err := t.st.Services.Projects.Refresh(t.logger)
		if err != nil {
			return doctorChecksLoadedMsg{err: errors.Wrap(err, "unable to load projects")}
		}
		checks.SetModules(t.st.Services.Modules.Deps(), t.st.Services.Modules.Dangerous())
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
		modules := t.projects.AllModules()
		if len(modules) == 0 {
			prjs, err := t.st.Services.Projects.Refresh(t.logger)
			if err != nil {
				return doctorAllResultsMsg{err: errors.Wrap(err, "unable to load projects")}
			}
			modules = prjs.AllModules()
		}
		checks.SetModules(t.st.Services.Modules.Deps(), t.st.Services.Modules.Dangerous())
		ret := checks.CheckAll(t.ctx, modules, t.logger)
		return doctorAllResultsMsg{results: ret}
	}
}

func projectDetails(p *project.Project) string {
	if p == nil {
		return "No project selected."
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Name: %s\n", p.Title()))
	b.WriteString(fmt.Sprintf("Key: %s\n", p.Key))
	b.WriteString(fmt.Sprintf("Package: %s\n", orDash(p.Package)))
	b.WriteString(fmt.Sprintf("Version: %s\n", orDash(p.Version)))
	b.WriteString(fmt.Sprintf("Exec: %s\n", orDash(p.ExecSafe())))
	b.WriteString(fmt.Sprintf("Path: %s\n", orDash(p.Path)))
	return b.String()
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
