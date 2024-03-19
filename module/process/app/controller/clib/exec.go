package clib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"{{{ if .HasModule "websocket" }}}
	"github.com/robert-nix/ansihtml"{{{ end }}}

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/exec"
	{{{ if .HasModule "websocket" }}}"{{{ .Package }}}/app/lib/websocket"
	{{{ end }}}"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vexec"
)

const execIcon = "file"

func ExecList(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Processes", as.Services.Exec.Execs)
		return controller.Render(w, r, as, &vexec.List{Execs: as.Services.Exec.Execs}, ps, "exec")
	})
}

func ExecForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		x := &exec.Exec{}
		ps.SetTitleAndData("New Process", x)
		ps.DefaultNavIcon = execIcon
		return controller.Render(w, r, as, &vexec.Form{Exec: x}, ps, "exec", "New Process")
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
		x := as.Services.Exec.NewExec(key, cmd, path, env...){{{ if .HasModule "websocket" }}}
		wf := func(key string, b []byte) error {
			m := util.ValueMap{"msg": string(b), "html": string(ansihtml.ConvertToHTML(b))}
			msg := &websocket.Message{Channel: key, Cmd: "output", Param: util.ToJSONBytes(m, true)}
			return as.Services.Socket.WriteChannel(msg, ps.Logger)
		}
		err = x.Start(wf){{{ else }}}
		err = x.Start(){{{ end }}}
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started process", x.WebPath(), w, ps)
	})
}

func ExecDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecRC(r, as)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ex.String(), ex)
		ps.DefaultNavIcon = execIcon
		return controller.Render(w, r, as, &vexec.Detail{Exec: ex}, ps, "exec", ex.String())
	})
}
{{{ if .HasModule "websocket" }}}
func ExecSocket(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ex, err := getExecRC(r, as)
		if err != nil {
			return "", err
		}
		err = as.Services.Socket.Upgrade(ps.Context, w, r, ex.String(){{{ if .HasUser }}}, ps.User{{{ end }}}, ps.Profile{{{ if .HasAccount }}}, ps.Accounts{{{ end }}}, ps.Logger)
		if err != nil {
			ps.Logger.Warn("unable to upgrade connection to websocket")
			return "", err
		}
		return "", nil
	})
}
{{{ end }}}
func ExecKill(w http.ResponseWriter, r *http.Request) {
	controller.Act("exec.kill", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		proc, err := getExecRC(r, as)
		if err != nil {
			return "", err
		}
		err = proc.Kill()
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, fmt.Sprintf("Killed process [%s]", proc.String()), "/admin/exec", w, ps)
	})
}

func getExecRC(r *http.Request, as *app.State) (*exec.Exec, error) {
	key, err := cutil.RCRequiredString(r, "key", false)
	if err != nil {
		return nil, err
	}
	idx, err := cutil.RCRequiredInt(r, "idx")
	if err != nil {
		return nil, err
	}
	proc := as.Services.Exec.Execs.Get(key, idx)
	if proc == nil {
		return nil, errors.Errorf("no process found with key [%s] and index [%d]", key, idx)
	}
	return proc, nil
}
