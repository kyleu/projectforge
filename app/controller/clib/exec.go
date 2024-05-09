// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexec"
)

const execIcon = "file"

func ExecList(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Processes", as.Services.Exec.Execs)
		return controller.Render(r, as, &vexec.List{Execs: as.Services.Exec.Execs}, ps, "exec")
	})
}

func ExecForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		x := &exec.Exec{}
		ps.SetTitleAndData("New Process", x)
		ps.DefaultNavIcon = execIcon
		return controller.Render(r, as, &vexec.Form{Exec: x}, ps, "exec", "New Process")
	})
}

func ExecNew(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.new", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		x := as.Services.Exec.NewExec(key, cmd, path, env...)
		err = x.Start(as.Services.Socket.Terminal(key, ps.Logger))
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started process", x.WebPath(), ps)
	})
}

func ExecDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecPath(r, as)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ex.String(), ex)
		ps.DefaultNavIcon = execIcon
		return controller.Render(r, as, &vexec.Detail{Exec: ex}, ps, "exec", ex.String())
	})
}

func ExecSocket(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecPath(r, as)
		if err != nil {
			return "", err
		}
		id, err := as.Services.Socket.Upgrade(ps.Context, w, r, ex.String(), ps.Profile, nil, ps.Logger)
		if err != nil {
			ps.Logger.Warn("unable to upgrade connection to websocket")
			return "", err
		}
		return "", as.Services.Socket.ReadLoop(ps.Context, id, ps.Logger)
	})
}

func ExecKill(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.kill", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		proc, err := getExecPath(r, as)
		if err != nil {
			return "", err
		}
		err = proc.Kill()
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, fmt.Sprintf("Killed process [%s]", proc.String()), "/admin/exec", ps)
	})
}

func getExecPath(r *http.Request, as *app.State) (*exec.Exec, error) {
	key, err := cutil.PathString(r, "key", false)
	if err != nil {
		return nil, err
	}
	idx, err := cutil.PathInt(r, "idx")
	if err != nil {
		return nil, err
	}
	proc := as.Services.Exec.Execs.Get(key, idx)
	if proc == nil {
		return nil, errors.Errorf("no process found with key [%s] and index [%d]", key, idx)
	}
	return proc, nil
}
