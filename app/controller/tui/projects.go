package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

var (
	screenProjects = NewScreen(
		"projects", "Projects", "", renderProjects,
		`"esc": back`, `"↑"/"↓": move`, `"enter": project details`, `"r": reload`, `"q": quit`,
	)
	screenProjectActions = NewScreen(
		"project-actions", "Project Actions", "", renderProjectActions,
		`"esc": back`, `"↑"/"↓": move`, `"enter": run action`, `"q": quit`,
	)
)

var projectActionTypes = []action.Type{action.TypePreview, action.TypeGenerate}

type projectsLoadedMsg struct {
	err error
}

type projectActionCompletedMsg struct {
	projectKey string
	actionKey  string
	result     *action.Result
	err        error
}

func renderProjects(t *TUI) string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Projects"))
	b.WriteString("\n\n")

	prjs := projectsFor(t)
	switch {
	case t.Config.projectsLoading:
		b.WriteString("Loading projects...")
	case t.Config.projectsErr != nil:
		b.WriteString("Error loading projects:\n")
		b.WriteString(t.Config.projectsErr.Error())
	case len(prjs) == 0:
		b.WriteString("No projects found.")
	default:
		items := make(menu.Items, 0, len(prjs))
		for _, p := range prjs {
			items = append(items, &menu.Item{Key: p.Key, Title: p.Title()})
		}
		cursor := t.Screen.Cursor()
		b.WriteString(RenderMenuOptions(cursor, items))
		b.WriteString("\n")
		if cursor >= 0 && cursor < len(prjs) {
			b.WriteString(resultStyle.Render(projectDetails(prjs[cursor])))
		}
	}

	return containerStyle.Width(t.width).Height(t.height).Render(b.String())
}

func renderProjectActions(t *TUI) string {
	var b strings.Builder
	prj := selectedProject(t)
	if prj == nil {
		b.WriteString(titleStyle.Render("Project"))
		b.WriteString("\n\nNo project selected.")
		return containerStyle.Width(t.width).Height(t.height).Render(b.String())
	}

	b.WriteString(titleStyle.Render(prj.Title()))
	b.WriteString("\n\n")

	items := make(menu.Items, 0, len(projectActionTypes))
	for _, a := range projectActionTypes {
		items = append(items, &menu.Item{Key: a.Key, Title: a.Title})
	}
	b.WriteString(RenderMenuOptions(t.Screen.Cursor(), items))
	b.WriteString("\n")

	cursor := t.Screen.Cursor()
	if cursor >= 0 && cursor < len(projectActionTypes) {
		a := projectActionTypes[cursor]
		b.WriteString(helpStyle.Render(a.Description))
		b.WriteString("\n\n")
	}

	if t.Config.projectActionRunning {
		b.WriteString(helpStyle.Render("Running action..."))
		b.WriteString("\n\n")
	}
	if t.Config.projectActionErr != nil {
		b.WriteString(resultStyle.Render("Action error:\n" + t.Config.projectActionErr.Error()))
		b.WriteString("\n\n")
	}

	if result := currentProjectActionResult(t, prj.Key, cursor); result != nil {
		b.WriteString(resultStyle.Render(projectActionDetails(result)))
		b.WriteString("\n\n")
	}

	b.WriteString(resultStyle.Render(projectDetails(prj)))
	return containerStyle.Width(t.width).Height(t.height).Render(b.String())
}

func onKeyProjects(key string, t *TUI) tea.Cmd {
	switch key {
	case "q":
		t.quitting = true
		return tea.Quit
	case tuiKeyEsc:
		t.Screen = screenMenu
	case tuiKeyUp, "k":
		if t.Screen.Cursor() > 0 {
			t.Screen.ModifyCursor(-1)
		}
	case tuiKeyDown, "j":
		prjs := projectsFor(t)
		if t.Screen.Cursor() < len(prjs)-1 {
			t.Screen.ModifyCursor(1)
		}
	case tuiKeyEnter, " ":
		prjs := projectsFor(t)
		cursor := t.Screen.Cursor()
		if cursor < 0 || cursor >= len(prjs) {
			return nil
		}
		t.Config.projectKey = prjs[cursor].Key
		t.Config.projectActionsFromProjects = true
		t.Config.projectActionErr = nil
		t.Config.projectActionRunning = false
		screenProjectActions.ResetCursor()
		t.Screen = screenProjectActions
	case "r":
		t.Config.projectsLoading = true
		t.Config.projectsErr = nil
		return loadProjectsCmd(t)
	}
	return nil
}

func onKeyProjectActions(key string, t *TUI) tea.Cmd {
	if key == "q" {
		t.quitting = true
		return tea.Quit
	}
	if key == tuiKeyEsc {
		t.Config.projectActionErr = nil
		if t.Config.projectActionsFromProjects {
			t.Screen = screenProjects
		} else {
			t.Screen = screenMenu
		}
		return nil
	}
	if t.Config.projectActionRunning {
		return nil
	}

	switch key {
	case tuiKeyUp, "k":
		if t.Screen.Cursor() > 0 {
			t.Screen.ModifyCursor(-1)
		}
	case tuiKeyDown, "j":
		if t.Screen.Cursor() < len(projectActionTypes)-1 {
			t.Screen.ModifyCursor(1)
		}
	case tuiKeyEnter, " ":
		prj := selectedProject(t)
		if prj == nil {
			return nil
		}
		cursor := t.Screen.Cursor()
		if cursor < 0 || cursor >= len(projectActionTypes) {
			return nil
		}
		a := projectActionTypes[cursor]
		t.Config.projectActionRunning = true
		t.Config.projectActionErr = nil
		return runProjectActionCmd(t, prj.Key, a)
	}
	return nil
}

func loadProjectsCmd(t *TUI) tea.Cmd {
	return func() tea.Msg {
		_, err := t.st.Services.Projects.Refresh(t.logger)
		return projectsLoadedMsg{err: err}
	}
}

func runProjectActionCmd(t *TUI, projectKey string, act action.Type) tea.Cmd {
	return func() tea.Msg {
		params := &action.Params{
			ProjectKey: projectKey,
			T:          act,
			Cfg:        util.ValueMap{},
			MSvc:       t.st.Services.Modules,
			PSvc:       t.st.Services.Projects,
			XSvc:       t.st.Services.Exec,
			SSvc:       t.st.Services.Socket,
			ESvc:       t.st.Services.Export,
			Logger:     t.logger,
		}
		res := action.Apply(t.ctx, params)
		return projectActionCompletedMsg{projectKey: projectKey, actionKey: act.Key, result: res}
	}
}

func projectsFor(t *TUI) project.Projects {
	return t.st.Services.Projects.Projects()
}

func selectedProject(t *TUI) *project.Project {
	if t.Config.projectKey == "" {
		return nil
	}
	for _, p := range projectsFor(t) {
		if p.Key == t.Config.projectKey {
			return p
		}
	}
	return nil
}

func currentProjectActionResult(t *TUI, projectKey string, cursor int) *action.Result {
	if t.Config.projectActionResults == nil {
		return nil
	}
	if cursor < 0 || cursor >= len(projectActionTypes) {
		return nil
	}
	key := projectActionResultKey(projectKey, projectActionTypes[cursor].Key)
	return t.Config.projectActionResults[key]
}

func projectActionResultKey(projectKey string, actionKey string) string {
	return projectKey + ":" + actionKey
}

func projectActionDetails(res *action.Result) string {
	if res == nil {
		return "No result yet. Press Enter to run this action."
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Action: %s\n", res.Action.Title))
	b.WriteString(fmt.Sprintf("Status: %s\n", strings.ToUpper(orDash(res.Status))))
	b.WriteString(fmt.Sprintf("Duration: %s\n", util.MicrosToMillis(res.Duration)))
	b.WriteString(fmt.Sprintf("Modules: %d\n", len(res.Modules)))
	if len(res.Errors) > 0 {
		b.WriteString("Errors:\n")
		for _, e := range res.Errors {
			b.WriteString("- ")
			b.WriteString(e)
			b.WriteByte('\n')
		}
	}
	if len(res.Logs) > 0 {
		b.WriteString("Logs:\n")
		for idx, l := range res.Logs {
			if idx >= 6 {
				b.WriteString("- ...\n")
				break
			}
			b.WriteString("- ")
			b.WriteString(l)
			b.WriteByte('\n')
		}
	}
	return strings.TrimSpace(b.String())
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
