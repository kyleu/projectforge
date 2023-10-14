// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vadmin"
)

func socketRoute(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState, path ...string) (string, error) {
	bc := func(extra ...string) []string {
		return append([]string{"admin", "Sockets||/admin/sockets"}, extra...)
	}
	if len(path) == 0 {
		ps.Title = "Sockets"
		chans, conns, taps := as.Services.Socket.Status()
		ps.Data = util.ValueMap{"channels": chans, "connections": conns}
		return controller.Render(rc, as, &vadmin.Sockets{Channels: chans, Connections: conns, Taps: taps}, ps, bc()...)
	}
	switch path[0] {
	case "tap":
		ps.Title = "WebSocket Tap"
		ps.Data = ps.Title
		return controller.Render(rc, as, &vadmin.SocketTap{}, ps, bc("Tap**search")...)
	case "tap-socket":
		_, err := as.Services.Socket.RegisterTap(rc, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to register tap socket")
		}
		return "", nil
	case "chan":
		if len(path) == 0 {
			return "", errors.New("no channel in path")
		}
		ch := as.Services.Socket.GetChannel(path[1])
		members := as.Services.Socket.GetAllMembers(path[1])
		ps.Data = ch
		return controller.Render(rc, as, &vadmin.Channel{Channel: ch, Members: members}, ps, bc("Channel", ch.Key+"**eye")...)
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
			ps.Data = c
			return controller.Render(rc, as, &vadmin.Connection{Connection: c}, ps, bc("Connection**eye", c.ID.String()+"**eye")...)
		}
		frm, _ := cutil.ParseForm(rc)
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
		return controller.FlashAndRedir(true, "sent message", fmt.Sprintf("/admin/sockets/conn/%s", id.String()), rc, ps)
	default:
		return "", errors.Errorf("invalid path [%s]", strings.Join(path, "/"))
	}
}
