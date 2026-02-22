package settings

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/util"
)

type execNewScreen struct {
	form  *huh.Form
	key   string
	cmd   string
	path  string
	env   string
	debug bool
}

func (s *execNewScreen) Key() string { return keyExecNew }

func (s *execNewScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "New Process"
	s.path = "."
	s.form = huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Key").Value(&s.key),
		huh.NewInput().Title("Command").Value(&s.cmd),
		huh.NewInput().Title("Path").Value(&s.path),
		huh.NewInput().Title("Env Vars (comma-separated)").Value(&s.env),
		huh.NewConfirm().Title("Enable debug output").Value(&s.debug),
	))
	ps.SetStatus("Fill the form and press enter")
	return s.form.Init()
}

func (s *execNewScreen) Update(ts *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	mdl, cmd := s.form.Update(msg)
	s.form = mdl.(*huh.Form)
	if s.form.State == huh.StateCompleted {
		if strings.TrimSpace(s.cmd) == "" {
			return mvc.Stay(), nil, errors.New("command is required")
		}
		key := strings.TrimSpace(s.key)
		if key == "" {
			key = strings.Fields(strings.TrimSpace(s.cmd))[0]
		}
		env := util.Str(s.env).SplitAndTrim(",").Strings()
		ex := ts.App.Services.Exec.NewExec(key, strings.TrimSpace(s.cmd), strings.TrimSpace(s.path), s.debug, env...)
		if err := ex.Start(ts.App.Services.Socket.Terminal(ex.String(), ts.Logger)); err != nil {
			return mvc.Stay(), nil, err
		}
		return mvc.Replace(keyExecDetail, util.ValueMap{"key": ex.String()}), nil, nil
	}
	return mvc.Stay(), cmd, nil
}

func (s *execNewScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	return renderPanel(style.New(ts.Theme), ps.Title, s.form.View(), rects)
}

func (s *execNewScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"tab: next", "enter: submit", "b: back"}}
}
