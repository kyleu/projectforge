package cproject

import (
	"fmt"
	"net/http"

	"github.com/robert-nix/ansihtml"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ProjectStart(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.start", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		cmd := fmt.Sprintf("%s%s -v server source:%s", project.DebugOutputDir, prj.Executable(), util.AppKey)
		exec := as.Services.Exec.NewExec(prj.Key, cmd, prj.Path, false)
		exec.Link = prj.WebPath()
		wf := func(key string, b []byte) error {
			m := util.ValueMap{"msg": string(b), "html": string(ansihtml.ConvertToHTML(b))}
			msg := &websocket.Message{Channel: key, Cmd: "output", Param: util.ToJSONBytes(m, true)}
			return as.Services.Socket.WriteChannel(msg, ps.Logger)
		}
		err = exec.Start(wf)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "started project", exec.WebPath(), ps)
	})
}
