// $PF_HAS_MODULE(process)$
package settings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/exec"
	"{{{ .Package }}}/app/util"
)

type execDetailScreen struct{ ex *exec.Exec }

func (s *execDetailScreen) Key() string { return keyExecDetail }

func (s *execDetailScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	s.ex = findExec(ts, ps.EnsureData().GetStringOpt("key"))
	if s.ex == nil {
		return func() tea.Msg { return errMsg{err: errors.New("unable to find process")} }
	}
	ps.Title = "Process " + s.ex.String()
	ps.SetStatus("k: kill process")
	return nil
}

func (s *execDetailScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "k" && s.ex != nil && s.ex.Completed == nil {
		if err := s.ex.Kill(); err != nil {
			return mvc.Stay(), nil, err
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *execDetailScreen) RefreshInterval(_ *mvc.State, _ *mvc.PageState) time.Duration {
	if s.ex == nil || s.ex.Completed != nil {
		return 0
	}
	return 200 * time.Millisecond
}

func (s *execDetailScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	if s.ex == nil {
		return renderPanel(st, ps.Title, "process not found", rects)
	}
	lines := []string{fmt.Sprintf("key: %s", s.ex.Key), fmt.Sprintf("pid: %d", s.ex.PID), "cmd: " + s.ex.Cmd, "path: " + s.ex.Path}
	if s.ex.Started != nil {
		lines = append(lines, "started: "+util.TimeRelative(s.ex.Started))
	}
	if s.ex.Completed != nil {
		lines = append(lines, "completed: "+util.TimeRelative(s.ex.Completed), fmt.Sprintf("exit: %d", s.ex.ExitCode))
	}
	lines = append(lines, "", "Output:", strings.TrimSpace(s.ex.Buffer.String()))
	return renderPanel(st, ps.Title, strings.Join(lines, "\n"), rects)
}

func (s *execDetailScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"k: kill", "b: back"}}
}

func findExec(ts *mvc.State, id string) *exec.Exec {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil
	}
	idx, _ := strconv.Atoi(parts[1])
	return ts.App.Services.Exec.Execs.Get(parts[0], idx)
}
