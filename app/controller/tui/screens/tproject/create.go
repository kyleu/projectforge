package tproject

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	pftheme "projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

type projectNewStep string

const (
	projectNewStepDir       = projectNewStep("dir")
	projectNewStepLoad      = projectNewStep("load")
	projectNewStepQuestions = projectNewStep("questions")
	projectNewStepCreate    = projectNewStep("create")
)

type projectNewLoadedMsg struct {
	dir string
	prj *project.Project
	err error
}

type projectNewCreatedMsg struct {
	projectKey string
	err        error
}

type ProjectNewScreen struct {
	step        projectNewStep
	form        *huh.Form
	dir         string
	prj         *project.Project
	prompts     []action.CreatePrompt
	promptIndex int
	textValue   string
	moduleValue []string
}

func NewProjectNewScreen() *ProjectNewScreen {
	return &ProjectNewScreen{step: projectNewStepDir}
}

func (s *ProjectNewScreen) Key() string {
	return KeyProjectNew
}

func (s *ProjectNewScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "New Project"
	ps.SetStatus("Choose a directory for your new project")
	s.step = projectNewStepDir
	s.dir = "."
	s.prj = nil
	s.prompts = nil
	s.promptIndex = 0
	s.textValue = ""
	s.moduleValue = nil
	s.form = s.newForm(ts, huh.NewGroup(
		huh.NewInput().
			Title("Directory").
			Description("Choose the directory your project will live in").
			Value(&s.dir),
	))
	return s.form.Init()
}

func (s *ProjectNewScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	switch m := msg.(type) {
	case projectNewLoadedMsg:
		if m.err != nil {
			return mvc.Stay(), nil, m.err
		}
		s.prj = m.prj
		s.dir = m.dir
		s.prompts = action.CreatePrompts()
		s.promptIndex = 0
		s.step = projectNewStepQuestions
		ps.SetStatus("Answer %d questions to configure your project", len(s.prompts))
		return mvc.Stay(), s.initPromptForm(ts), nil
	case projectNewCreatedMsg:
		if m.err != nil {
			s.step = projectNewStepCreate
			ps.SetStatus("Creation failed")
			return mvc.Stay(), nil, m.err
		}
		ps.SetStatus("Project created")
		return mvc.Replace(KeyProject, util.ValueMap{"project": m.projectKey}), nil, nil
	case tea.KeyMsg:
		if m.String() == screens.KeyEsc {
			if s.step == projectNewStepCreate || s.step == projectNewStepLoad {
				return mvc.Stay(), nil, nil
			}
			if s.step == projectNewStepQuestions && s.promptIndex > 0 {
				s.promptIndex--
				return mvc.Stay(), s.initPromptForm(ts), nil
			}
			if s.step == projectNewStepQuestions && s.promptIndex == 0 {
				s.step = projectNewStepDir
				ps.SetStatus("Choose a directory for your new project")
				s.form = s.newForm(ts, huh.NewGroup(
					huh.NewInput().
						Title("Directory").
						Description("Choose the directory your project will live in").
						Value(&s.dir),
				))
				return mvc.Stay(), s.form.Init(), nil
			}
			return mvc.Pop(), nil, nil
		}
	}
	if s.form == nil {
		return mvc.Stay(), nil, nil
	}
	mdl, cmd := s.form.Update(msg)
	form, ok := mdl.(*huh.Form)
	if !ok {
		return mvc.Stay(), nil, errors.Errorf("invalid [%T] as model, expected [*Form]", mdl)
	}
	s.form = form
	if s.form.State != huh.StateCompleted {
		return mvc.Stay(), cmd, nil
	}
	switch s.step {
	case projectNewStepDir:
		s.dir = util.OrDefault(strings.TrimSpace(s.dir), ".")
		ps.SetStatus("Loading project from [%s]", s.dir)
		s.step = projectNewStepLoad
		return mvc.Stay(), s.loadProjectCmd(ts, ps), nil
	case projectNewStepQuestions:
		if err := s.applyCurrentPrompt(); err != nil {
			return mvc.Stay(), nil, err
		}
		s.promptIndex++
		if s.promptIndex >= len(s.prompts) {
			s.step = projectNewStepCreate
			ps.SetStatus("Creating and building project...")
			return mvc.Stay(), s.createProjectCmd(ts, ps), nil
		}
		return mvc.Stay(), s.initPromptForm(ts), nil
	case projectNewStepLoad, projectNewStepCreate:
		return mvc.Stay(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectNewScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	var body strings.Builder
	switch s.step {
	case projectNewStepLoad:
		body.WriteString("Loading project information...\n")
	case projectNewStepCreate:
		body.WriteString("Creating project and running full build...\n")
	default:
		if s.form != nil {
			body.WriteString(s.form.View())
			body.WriteString("\n")
		}
	}
	showSummary := s.prj != nil
	if s.step == projectNewStepQuestions && s.currentPromptKind() == action.CreatePromptModules {
		// Keep module picker usable in smaller terminals by giving it full panel height.
		showSummary = false
	}
	if showSummary {
		body.WriteString("\nSummary:\n")
		body.WriteString(fmt.Sprintf("Directory: %s\n", s.dir))
		body.WriteString(strings.Join(action.CreateSummaryLines(s.prj), "\n"))
	}
	return renderScreenPanel(ps.Title, body.String(), styles.Panel, styles, rects)
}

func (s *ProjectNewScreen) Help(_ *mvc.State, _ *mvc.PageState) HelpBindings {
	return HelpBindings{Short: []string{"tab: next", "enter: submit", "esc: back"}}
}

func (s *ProjectNewScreen) BypassGlobalKey(_ *mvc.State, _ *mvc.PageState, key string) bool {
	return key == "/" && s.form != nil
}

func (s *ProjectNewScreen) loadProjectCmd(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	return func() tea.Msg {
		logger := ps.Logger
		if logger == nil {
			logger = ts.Logger
		}
		_, _ = ts.App.Services.Projects.Refresh(logger)
		prj := ts.App.Services.Projects.ByPath(s.dir)
		if prj == nil {
			prj = project.NewProject(filepath.Base(s.dir), s.dir)
		}
		prj.Name = strings.TrimSuffix(prj.Name, " (missing)")
		return projectNewLoadedMsg{dir: s.dir, prj: prj}
	}
}

func (s *ProjectNewScreen) createProjectCmd(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	return func() tea.Msg {
		if s.prj == nil {
			return projectNewCreatedMsg{err: errors.New("project is not initialized")}
		}
		logger := ps.Logger
		if logger == nil {
			logger = ts.Logger
		}
		cfg := action.CreateConfigFromProject(s.prj)
		params := &action.Params{
			ProjectKey: s.prj.Key,
			T:          action.TypeCreate,
			Cfg:        cfg,
			MSvc:       ts.App.Services.Modules,
			PSvc:       ts.App.Services.Projects,
			ESvc:       ts.App.Services.Export,
			XSvc:       ts.App.Services.Exec,
			Logger:     logger,
		}
		res := action.Apply(ps.Context, params)
		if res == nil {
			return projectNewCreatedMsg{err: errors.New("create returned no result")}
		}
		if err := res.AsError(); err != nil {
			return projectNewCreatedMsg{err: err}
		}
		createdKey := ""
		if res.Project != nil {
			createdKey = res.Project.Key
		}
		if createdKey == "" {
			createdKey = projectKeyByPath(ts, s.dir)
		}
		if createdKey == "" {
			createdKey = s.prj.Key
		}
		if createdKey == "" {
			return projectNewCreatedMsg{err: errors.New("project created without a key")}
		}
		return projectNewCreatedMsg{projectKey: createdKey}
	}
}

func (s *ProjectNewScreen) initPromptForm(ts *mvc.State) tea.Cmd {
	if s.prj == nil || s.promptIndex < 0 || s.promptIndex >= len(s.prompts) {
		return nil
	}
	q := s.prompts[s.promptIndex]
	switch q.Kind {
	case action.CreatePromptModules:
		s.moduleValue = action.CreatePromptDefaultModules(s.prj)
		opts := make([]huh.Option[string], 0, len(ts.App.Services.Modules.Modules()))
		for _, mod := range ts.App.Services.Modules.Modules().Sorted() {
			opts = append(opts, huh.NewOption(fmt.Sprintf("%s: %s", mod.Key, mod.Description), mod.Key))
		}
		s.form = s.newForm(ts, huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(fmt.Sprintf("%d: %s", s.promptIndex+1, q.Query)).
				Description("Use arrow keys to move and space to toggle modules").
				Height(12).
				Options(opts...).
				Value(&s.moduleValue),
		))
	default:
		s.textValue = action.CreatePromptDefaultString(s.prj, q.Key)
		in := huh.NewInput().Title(fmt.Sprintf("%d: %s", s.promptIndex+1, q.Query)).Value(&s.textValue)
		if s.textValue == "" {
			in = in.Description("Optional")
		} else {
			in = in.Description("Default: " + s.textValue)
		}
		s.form = s.newForm(ts, huh.NewGroup(in))
	}
	return s.form.Init()
}

func (s *ProjectNewScreen) applyCurrentPrompt() error {
	if s.prj == nil || s.promptIndex < 0 || s.promptIndex >= len(s.prompts) {
		return nil
	}
	q := s.prompts[s.promptIndex]
	switch q.Kind {
	case action.CreatePromptModules:
		action.ApplyCreatePromptModules(s.prj, s.moduleValue)
	default:
		val := util.OrDefault(strings.TrimSpace(s.textValue), action.CreatePromptDefaultString(s.prj, q.Key))
		if err := action.ApplyCreatePromptString(s.prj, q.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProjectNewScreen) currentPromptKind() action.CreatePromptKind {
	if s.promptIndex < 0 || s.promptIndex >= len(s.prompts) {
		return action.CreatePromptString
	}
	return s.prompts[s.promptIndex].Kind
}

func (s *ProjectNewScreen) newForm(ts *mvc.State, group *huh.Group) *huh.Form {
	return huh.NewForm(group).WithTheme(projectNewHuhTheme(ts))
}

func projectNewHuhTheme(ts *mvc.State) *huh.Theme {
	c := projectNewColors(ts)
	bg := lipgloss.Color(c.Background)
	bgm := lipgloss.Color(c.BackgroundMuted)
	fg := lipgloss.Color(c.Foreground)
	fgm := lipgloss.Color(c.ForegroundMuted)
	nbg := lipgloss.Color(c.NavBackground)
	nfg := lipgloss.Color(c.NavForeground)

	t := huh.ThemeCharm()
	t.Form.Base = t.Form.Base.Background(bg).Foreground(fg)
	t.Group.Base = t.Group.Base.Background(bg).Foreground(fg)
	t.Group.Title = t.Group.Title.Foreground(nfg).Background(bg)
	t.Group.Description = t.Group.Description.Foreground(fgm).Background(bg)

	applyField := func(fs *huh.FieldStyles) {
		fs.Base = fs.Base.Background(bg).Foreground(fg)
		fs.Title = fs.Title.Background(bg)
		fs.Description = fs.Description.Foreground(fgm).Background(bg)
		fs.Option = fs.Option.Background(bg)
		fs.SelectedOption = fs.SelectedOption.Background(bg)
		fs.UnselectedOption = fs.UnselectedOption.Background(bg)
		fs.SelectedPrefix = fs.SelectedPrefix.Background(bg)
		fs.UnselectedPrefix = fs.UnselectedPrefix.Background(bg)
		fs.FocusedButton = fs.FocusedButton.Foreground(nfg).Background(nbg)
		fs.BlurredButton = fs.BlurredButton.Foreground(fg).Background(bgm)
		fs.Card = fs.Card.Background(bg).Foreground(fg)
		fs.TextInput.Text = fs.TextInput.Text.Foreground(fg).Background(bg)
		fs.TextInput.Prompt = fs.TextInput.Prompt.Background(bg)
		fs.TextInput.Placeholder = fs.TextInput.Placeholder.Background(bg)
	}
	applyField(&t.Focused)
	applyField(&t.Blurred)
	return t
}

func projectNewColors(ts *mvc.State) *pftheme.Colors {
	th := pftheme.Default
	if ts != nil && ts.Theme != nil {
		th = ts.Theme
	}
	dark := style.IsDarkMode()
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
	return pftheme.Default.Dark
}

func projectKeyByPath(ts *mvc.State, path string) string {
	if ts == nil || ts.App == nil || ts.App.Services == nil || ts.App.Services.Projects == nil {
		return ""
	}
	currAbs, currErr := filepath.Abs(path)
	for _, p := range ts.App.Services.Projects.Projects() {
		if p == nil {
			continue
		}
		if p.Path == path {
			return p.Key
		}
		if currErr == nil {
			if pAbs, err := filepath.Abs(p.Path); err == nil && pAbs == currAbs {
				return p.Key
			}
		}
	}
	return ""
}
