package cproject

import (
	"fmt"

	"github.com/robert-nix/ansihtml"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
)

func ProjectStart(rc *fasthttp.RequestCtx) {
	controller.Act("project.start", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		cmd := fmt.Sprintf("build/debug/%s -v server source:%s", prj.Executable(), util.AppKey)
		exec, err := as.Services.Exec.NewExec(prj.Key, cmd, prj.Path)
		if err != nil {
			return "", err
		}
		exec.Link = "/p/" + prj.Key
		w := func(key string, b []byte) error {
			m := util.ValueMap{"msg": string(b), "html": string(ansihtml.ConvertToHTML(b))}
			msg := &websocket.Message{Channel: key, Cmd: "output", Param: util.ToJSONBytes(m, true)}
			return as.Services.Socket.WriteChannel(msg)
		}
		err = exec.Start(ps.Context, ps.Logger, w)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started project", exec.WebPath(), rc, ps)
	})
}
