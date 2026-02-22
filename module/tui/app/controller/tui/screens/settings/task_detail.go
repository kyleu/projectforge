// $PF_HAS_MODULE(task)$
package settings

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/task"
	"{{{ .Package }}}/app/util"
)

type taskDetailScreen struct{ task *task.Task }

func (s *taskDetailScreen) Key() string { return keyTaskDetail }

func (s *taskDetailScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	s.task = ts.App.Services.Task.RegisteredTasks.Get(ps.EnsureData().GetStringOpt("key"))
	if s.task == nil {
		return func() tea.Msg { return errMsg{err: errors.New("unable to find task")} }
	}
	ps.Title = s.task.TitleSafe()
	ps.SetStatus("enter: run task")
	return nil
}

func (s *taskDetailScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "enter" && s.task != nil {
		return mvc.Push(keyTaskRun, util.ValueMap{"key": s.task.Key}), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *taskDetailScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	if s.task == nil {
		return renderPanel(st, ps.Title, "task not found", rects)
	}
	lines := []string{"key: " + s.task.Key, "category: " + s.task.Category, "description: " + s.task.Description, ""}
	for _, f := range s.task.Fields {
		lines = append(lines, fmt.Sprintf("arg %s (%s) default=%s", f.Key, f.Type, f.Default))
	}
	return renderPanel(st, ps.Title, strings.Join(lines, "\n"), rects)
}

func (s *taskDetailScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"enter: run", "b: back"}}
}
