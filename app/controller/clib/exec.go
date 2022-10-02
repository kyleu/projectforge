// Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/exec"
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
		ps.Title = "Processes"
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
		err = x.Start(ps.Context, nil, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started process", x.WebPath(), rc, ps)
	})
}

func ExecDetail(rc *fasthttp.RequestCtx) {
	controller.Act("exec.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		idx, err := cutil.RCRequiredInt(rc, "idx")
		if err != nil {
			return "", err
		}
		proc := as.Services.Exec.Execs.Get(key, idx)
		if proc == nil {
			return "", errors.Errorf("no process found with key [%s] and index [%d]", key, idx)
		}
		ps.Title = proc.String()
		ps.Data = proc
		return controller.Render(rc, as, &vexec.Detail{Exec: proc}, ps, "exec", proc.String())
	})
}
