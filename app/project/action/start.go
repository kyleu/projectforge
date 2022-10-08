package action

import (
	"context"
	"fmt"

	"github.com/robert-nix/ansihtml"

	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
)

func onStart(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
	if ret == nil {
		ret = newResult(TypeDebug, pm.Prj, pm.Cfg, pm.Logger)
	}

	cmd := fmt.Sprintf("build/debug/%s -v server source:%s", pm.Prj.Executable(), util.AppKey)
	exec, err := pm.XSvc.NewExec(pm.Prj.Key, cmd, pm.Prj.Path, pm.Prj.Info.EnvVars...)
	if err != nil {
		return errorResult(err, TypeDebug, pm.Cfg, pm.Logger)
	}
	exec.Link = "/p/" + pm.Prj.Key
	w := func(key string, b []byte) error {
		m := util.ValueMap{"msg": string(b), "html": string(ansihtml.ConvertToHTML(b))}
		msg := &websocket.Message{Channel: key, Cmd: "output", Param: util.ToJSONBytes(m, true)}
		if pm.SSvc == nil {
			return nil
		}
		return pm.SSvc.WriteChannel(msg)
	}
	err = exec.Start(ctx, pm.Logger, w)
	if err != nil {
		return errorResult(err, TypeDebug, pm.Cfg, pm.Logger)
	}

	return ret
}
