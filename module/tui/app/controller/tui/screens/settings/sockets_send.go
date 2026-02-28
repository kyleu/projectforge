// $PF_HAS_MODULE(websocket)$
package settings

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/websocket"
	"{{{ .Package }}}/app/util"
)

type socketSendScreen struct {
	form    *huh.Form
	to      string
	from    string
	channel string
	cmd     string
	param   string
}

func (s *socketSendScreen) Key() string { return keySocketSend }

func (s *socketSendScreen) Init(_ *mvc.State, ps *mvc.PageState) tea.Cmd {
	s.to = ps.EnsureData().GetStringOpt("id")
	s.from = s.to
	s.form = huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("To Connection ID").Value(&s.to),
		huh.NewInput().Title("From Connection ID").Value(&s.from),
		huh.NewInput().Title("Channel").Value(&s.channel),
		huh.NewInput().Title("Command").Value(&s.cmd),
		huh.NewInput().Title("Param JSON").Value(&s.param),
	))
	ps.Title = "Send Message"
	ps.SetStatus("Submit form to send")
	return s.form.Init()
}

func (s *socketSendScreen) Update(ts *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == screens.KeyEsc || m.String() == screens.KeyBackspace || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	mdl, cmd := s.form.Update(msg)
	form, ok := mdl.(*huh.Form)
	if !ok {
		return mvc.Stay(), nil, errors.New("invalid form model in [socket]")
	}
	s.form = form
	if s.form.State == huh.StateCompleted {
		to := util.UUIDFromString(s.to)
		if to == nil {
			return mvc.Stay(), nil, errors.Errorf("invalid destination connection id [%s]", s.to)
		}
		from := util.UUIDFromString(s.from)
		m := &websocket.Message{From: from, Channel: s.channel, Cmd: s.cmd, Param: []byte(s.param)}
		if err := ts.App.Services.Socket.WriteMessage(*to, m, ts.Logger); err != nil {
			return mvc.Stay(), nil, err
		}
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), cmd, nil
}

func (s *socketSendScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	return renderPanel(style.New(ts.Theme), ps.Title, s.form.View(), rects)
}

func (s *socketSendScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"tab: next", "enter: submit", "b: back"}}
}
