package clib

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vadmin"
)

const socketIcon = "eye"

func socketRoute(w http.ResponseWriter, r *http.Request, as *app.State, ps *cutil.PageState, path ...string) (string, error) {
	bc := func(extra ...string) []string {
		return append([]string{keyAdmin, "Sockets||/admin/sockets"}, extra...)
	}
	if len(path) == 0 {
		chans, conns, taps := as.Services.Socket.Status()
		ps.SetTitleAndData("Sockets", util.ValueMap{"channels": chans, "connections": conns})
		return controller.Render(r, as, &vadmin.Sockets{Channels: chans, Connections: conns, Taps: taps}, ps, bc()...)
	}
	switch path[0] {
	case "tap":
		ps.SetTitleAndData("WebSocket Tap", "tap")
		ps.DefaultNavIcon = "search"
		return controller.Render(r, as, &vadmin.SocketTap{}, ps, bc("Tap")...)
	case "tap-socket":
		_, err := as.Services.Socket.RegisterTap(w, r, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to register tap socket")
		}
		return "", nil
	case "chan":
		if len(path) == 0 {
			return "", errors.New("no channel in path")
		}
		ch := as.Services.Socket.GetChannel(path[1])
		if ch == nil {
			return "", errors.Errorf("can't find channel [%s]", path[1])
		}
		members := as.Services.Socket.GetAllMembers(ch.Key)
		ps.SetTitleAndData("Channel ["+ch.Key+"]", ch)
		ps.DefaultNavIcon = socketIcon
		return controller.Render(r, as, &vadmin.Channel{Channel: ch, Members: members}, ps, bc("Channel", ch.Key)...)
	case "conn":
		if len(path) == 0 {
			return "", errors.New("no connection ID in path")
		}
		id := util.UUIDFromString(path[1])
		if id == nil {
			return "", errors.Errorf("invalid connection ID [%s] in path", path[1])
		}
		if len(path) == 2 {
			c := as.Services.Socket.GetConnection(*id)
			if c == nil {
				return "", errors.Errorf("no connection with ID [%s]", id.String())
			}
			ps.SetTitleAndData("Connection ["+c.ID.String()+"]", c)
			ps.DefaultNavIcon = socketIcon
			return controller.Render(r, as, &vadmin.Connection{Connection: c}, ps, bc("Connection", c.ID.String())...)
		}
		frm, _ := cutil.ParseForm(r, ps.RequestBody)
		fromStr := frm.GetStringOpt("from")
		from := util.UUIDFromString(fromStr)
		channel := frm.GetStringOpt("channel")
		cmd := frm.GetStringOpt("cmd")
		param := frm.GetStringOpt("param")
		m := &websocket.Message{
			From:    from,
			Channel: channel,
			Cmd:     cmd,
			Param:   []byte(param),
		}
		err := as.Services.Socket.WriteMessage(*id, m, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "sent message", fmt.Sprintf("/admin/sockets/conn/%s", id.String()), ps)
	default:
		return "", errors.Errorf("invalid path [%s]", util.StringJoin(path, "/"))
	}
}
