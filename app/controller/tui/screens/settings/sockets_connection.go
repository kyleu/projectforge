package settings

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
)

type socketConnectionScreen struct{ conn *websocket.Connection }

func (s *socketConnectionScreen) Key() string { return keySocketConn }

func (s *socketConnectionScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	id := util.UUIDFromString(ps.EnsureData().GetStringOpt("key"))
	if id != nil {
		s.conn = ts.App.Services.Socket.GetConnection(*id)
	}
	if s.conn == nil {
		return func() tea.Msg { return errMsg{err: errors.New("connection not found")} }
	}
	ps.Title = "Connection " + s.conn.ID.String()
	ps.SetStatus("s: send message")
	return nil
}

func (s *socketConnectionScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if em, ok := msg.(errMsg); ok {
		return mvc.Stay(), nil, em.err
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "s" && s.conn != nil {
		return mvc.Push(keySocketSend, util.ValueMap{"id": s.conn.ID.String()}), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *socketConnectionScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	if s.conn == nil {
		return renderPanel(style.New(ts.Theme), ps.Title, "connection not found", rects)
	}
	lines := []string{fmt.Sprintf("id: %s", s.conn.ID.String()), "user: " + s.conn.Username(), "service: " + s.conn.Svc, "channels: " + strings.Join(s.conn.Channels, ", ")}
	if s.conn.Stats != nil {
		lines = append(lines, fmt.Sprintf("sent: %d msgs (%s)", s.conn.Stats.MessagesSent, util.ByteSizeSI(int64(s.conn.Stats.BytesSent))))
		lines = append(lines, fmt.Sprintf("recv: %d msgs (%s)", s.conn.Stats.MessagesReceived, util.ByteSizeSI(int64(s.conn.Stats.BytesReceived))))
	}
	return renderPanel(style.New(ts.Theme), ps.Title, strings.Join(lines, "\n"), rects)
}

func (s *socketConnectionScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"s: send", "b: back"}}
}
