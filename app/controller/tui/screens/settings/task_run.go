package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/app/util"
)

type taskRunMsg struct {
	result *task.Result
	err    error
}

type taskRunScreen struct {
	form     *huh.Form
	task     *task.Task
	category string
	values   map[string]*string
	result   *task.Result
}

func (s *taskRunScreen) Key() string { return keyTaskRun }

func (s *taskRunScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	key := ps.EnsureData().GetStringOpt("key")
	s.task = ts.App.Services.Task.RegisteredTasks.Get(key)
	if s.task == nil {
		return func() tea.Msg { return errMsg{err: errors.Errorf("unable to find task [%s]", key)} }
	}
	ps.Title = "Run " + s.task.TitleSafe()
	s.category = "ad-hoc"
	s.values = map[string]*string{}
	fields := []huh.Field{huh.NewInput().Title("Category").Value(&s.category)}
	for _, fd := range s.task.Fields {
		v := fd.Default
		s.values[fd.Key] = &v
		fields = append(fields, huh.NewInput().Title(fd.TitleSafe()).Description(fd.Description).Value(&v))
	}
	s.form = huh.NewForm(huh.NewGroup(fields...))
	ps.SetStatus("Submit form to run")
	return s.form.Init()
}

func (s *taskRunScreen) Update(ts *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if res, ok := msg.(taskRunMsg); ok {
		s.result = res.result
		return mvc.Stay(), nil, res.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == screens.KeyEsc || m.String() == screens.KeyBackspace || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	mdl, cmd := s.form.Update(msg)
	form, ok := mdl.(*huh.Form)
	if !ok {
		return mvc.Stay(), nil, errors.New("invalid form model in [task]")
	}
	s.form = form
	if s.form.State == huh.StateCompleted {
		args := util.ValueMap{}
		for k, v := range s.values {
			args[k] = strings.TrimSpace(*v)
		}
		return mvc.Stay(), func() tea.Msg { return taskRunMsg{result: s.task.Run(ts.Context, s.category, args, ts.Logger)} }, nil
	}
	return mvc.Stay(), cmd, nil
}

func (s *taskRunScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	body := s.form.View()
	if s.result != nil {
		body += "\n\nResult:\n"
		body += fmt.Sprintf("status: %s\nsummary: %s\n", s.result.Status, s.result.Summarize())
		body += strings.Join(s.result.Logs, "\n")
	}
	return renderPanel(style.New(ts.Theme), ps.Title, body, rects)
}

func (s *taskRunScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"tab: next", "enter: submit", "b: back"}}
}
