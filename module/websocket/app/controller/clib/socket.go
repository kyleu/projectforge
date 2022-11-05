package clib

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/websocket"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vadmin"
)

func socketRoute(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState, path ...string) (string, error) {
	bc := func(extra ...string) []string {
		return append([]string{"admin", "Sockets||/admin/sockets"}, extra...)
	}
	if len(path) == 0 {
		ps.Title = "Sockets"
		chans, conns, ctx := as.Services.Socket.Status()
		ps.Data = util.ValueMap{"channels": chans, "connections": conns, "context": ctx}
		return controller.Render(rc, as, &vadmin.Sockets{Channels: chans, Connections: conns, Context: ctx}, ps, bc()...)
	}
	switch path[0] {
	case "tap":
		return controller.Render(rc, as, &vadmin.Tap{}, ps, bc("Tap")...)
	case "chan":
		if len(path) == 0 {
			return "", errors.New("no channel in path")
		}
		ch := as.Services.Socket.GetChannel(path[1])
		members := as.Services.Socket.GetAllMembers(path[1])
		ps.Data = ch
		return controller.Render(rc, as, &vadmin.Channel{Channel: ch, Members: members}, ps, bc("Channel", ch.Key)...)
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
			ps.Data = c
			return controller.Render(rc, as, &vadmin.Connection{Connection: c}, ps, bc("Connection", c.ID.String())...)
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
		err := as.Services.Socket.WriteMessage(*id, m)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "sent message", fmt.Sprintf("/admin/sockets/conn/%s", id.String()), rc, ps)
	default:
		return "", errors.Errorf("invalid path [%s]", strings.Join(path, "/"))
	}
}
