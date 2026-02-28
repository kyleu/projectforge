package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/screens"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
)

const (
	keySocketChannel = "settings.admin.sockets.channel"
	keySocketConn    = "settings.admin.sockets.connection"
	keySocketSend    = "settings.admin.sockets.send"
)

type socketListScreen struct{ items menu.Items }

func registerSocketScreens(reg *screens.Registry) {
	reg.AddScreen(&socketListScreen{})
	reg.AddScreen(&socketChannelScreen{})
	reg.AddScreen(&socketConnectionScreen{})
	reg.AddScreen(&socketSendScreen{})
}

func (s *socketListScreen) Key() string { return keySockets }

func (s *socketListScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = "Active WebSockets"
	chans, conns, taps := ts.App.Services.Socket.Status()
	s.items = menu.Items{{Key: "taps", Title: fmt.Sprintf("Activity Taps (%d)", len(taps)), Description: "Read-only tap list", Route: ""}}
	for _, ch := range chans {
		s.items = append(s.items, &menu.Item{Key: ch, Title: "Channel: " + ch, Description: "View channel members", Route: keySocketChannel})
	}
	for _, c := range conns {
		s.items = append(s.items, &menu.Item{Key: c.ID.String(), Title: "Connection: " + c.ID.String(), Description: c.Username(), Route: keySocketConn})
	}
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	ps.SetStatus("channels=%d connections=%d taps=%d", len(chans), len(conns), len(taps))
	return nil
}

func (s *socketListScreen) Update(_ *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if delta, ok := menuDelta(msg); ok {
		ps.Cursor = menuClamp(ps.Cursor+delta, len(s.items))
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == screens.KeyEnter && len(s.items) > 0 {
		item := s.items[menuClamp(ps.Cursor, len(s.items))]
		if item.Route != "" {
			return mvc.Push(item.Route, util.ValueMap{"key": item.Key}), nil, nil
		}
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == screens.KeyEsc || m.String() == screens.KeyBackspace || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *socketListScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	return renderPanel(st, ps.Title, renderMenuBody(s.items, ps.Cursor, st, rects), rects)
}

func (s *socketListScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"up/down: move", "enter: open", "b: back"}}
}
