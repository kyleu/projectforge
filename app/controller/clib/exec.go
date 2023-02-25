// Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/robert-nix/ansihtml"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexec"
)

func ExecList(rc *fasthttp.RequestCtx) {
	controller.Act("exec.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Processes"
		ps.Data = as.Services.Exec.Execs
		return controller.Render(rc, as, &vexec.List{Execs: as.Services.Exec.Execs}, ps, "exec")
	})
}

func ExecForm(rc *fasthttp.RequestCtx) {
	controller.Act("exec.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		x := &exec.Exec{}
		ps.Title = "New Process"
		ps.Data = x
		return controller.Render(rc, as, &vexec.Form{Exec: x}, ps, "exec", "New Process")
	})
}

func ExecNew(rc *fasthttp.RequestCtx) {
	controller.Act("exec.new", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		cmd := strings.TrimSpace(frm.GetStringOpt("cmd"))
		if cmd == "" {
			return "", errors.New("must provide non-empty [cmd]")
		}
		key := strings.TrimSpace(frm.GetStringOpt("key"))
		if key == "" {
			key, _ = util.StringSplit(cmd, ' ', true)
		}
		path := strings.TrimSpace(frm.GetStringOpt("path"))
		if path == "" {
			path = "."
		}
		env := util.StringSplitAndTrim(strings.TrimSpace(frm.GetStringOpt("env")), ",")
		x, err := as.Services.Exec.NewExec(key, cmd, path, env...)
		if err != nil {
			return "", err
		}
		w := func(key string, b []byte) error {
			m := util.ValueMap{"msg": string(b), "html": string(ansihtml.ConvertToHTML(b))}
			msg := &websocket.Message{Channel: key, Cmd: "output", Param: util.ToJSONBytes(m, true)}
			return as.Services.Socket.WriteChannel(msg, ps.Logger)
		}
		err = x.Start(ps.Context, ps.Logger, w)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started process", x.WebPath(), rc, ps)
	})
}

func ExecDetail(rc *fasthttp.RequestCtx) {
	controller.Act("exec.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecRC(rc, as)
		if err != nil {
			return "", err
		}
		ps.Title = ex.String()
		ps.Data = ex
		return controller.Render(rc, as, &vexec.Detail{Exec: ex}, ps, "exec", ex.String())
	})
}

func ExecSocket(rc *fasthttp.RequestCtx) {
	controller.Act("exec.socket", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecRC(rc, as)
		if err != nil {
			return "", err
		}
		err = as.Services.Socket.Upgrade(ps.Context, rc, ex.String(), ps.Profile, ps.Logger)
		if err != nil {
			ps.Logger.Warn("unable to upgrade connection to websocket")
			return "", err
		}
		return "", nil
	})
}

func ExecKill(rc *fasthttp.RequestCtx) {
	controller.Act("exec.kill", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		proc, err := getExecRC(rc, as)
		if err != nil {
			return "", err
		}
		err = proc.Kill()
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, fmt.Sprintf("Killed process [%s]", proc.String()), "/admin/exec", rc, ps)
	})
}

func getExecRC(rc *fasthttp.RequestCtx, as *app.State) (*exec.Exec, error) {
	key, err := cutil.RCRequiredString(rc, "key", false)
	if err != nil {
		return nil, err
	}
	idx, err := cutil.RCRequiredInt(rc, "idx")
	if err != nil {
		return nil, err
	}
	proc := as.Services.Exec.Execs.Get(key, idx)
	if proc == nil {
		return nil, errors.Errorf("no process found with key [%s] and index [%d]", key, idx)
	}
	return proc, nil
}
