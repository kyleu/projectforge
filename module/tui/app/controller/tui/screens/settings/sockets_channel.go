// $PF_HAS_MODULE(websocket)$
package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
)

type socketChannelScreen struct {
	key   string
	lines []string
}

func (s *socketChannelScreen) Key() string { return keySocketChannel }

func (s *socketChannelScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	s.key = ps.EnsureData().GetStringOpt("key")
	ch := ts.App.Services.Socket.GetChannel(s.key)
	ps.Title = "Channel " + s.key
	s.lines = []string{"channel: " + s.key}
	for _, m := range ts.App.Services.Socket.GetAllMembers(s.key) {
		s.lines = append(s.lines, "member: "+m.ID.String()+" ("+m.Username()+")")
	}
	if ch != nil {
		ids := make([]string, 0, len(ch.ConnIDs))
		for _, id := range ch.ConnIDs {
			ids = append(ids, id.String())
		}
		s.lines = append(s.lines, fmt.Sprintf("connections: %s", strings.Join(ids, ", ")))
	}
	ps.SetStatus("Channel details")
	return nil
}

func (s *socketChannelScreen) Update(_ *mvc.State, _ *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *socketChannelScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	return renderPanel(style.New(ts.Theme), ps.Title, strings.Join(s.lines, "\n"), rects)
}

func (s *socketChannelScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"b: back"}}
}
