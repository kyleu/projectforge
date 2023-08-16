package clib

import (
	"strings"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vscripting"
)

func ScriptingList(rc *fasthttp.RequestCtx) {
	controller.Act("scripting.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Scripting"
		ret := as.Services.Script.ListScripts(ps.Logger)
		ps.Data = ret
		return controller.Render(rc, as, &vscripting.List{Scripts: ret}, ps, "scripting")
	})
}

var Examples = [][]any{{"a"}, {"b"}, {"c"}}

func ScriptingDetail(rc *fasthttp.RequestCtx) {
	controller.Act("scripting.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", true)
		if err != nil {
			return "", err
		}
		src, err := as.Services.Script.LoadScript(key, ps.Logger)
		if err != nil {
			return "", err
		}
		res := make(map[string]any, len(Examples))
		for _, ex := range Examples {
			x, err := as.Services.Script.RunScript(src, "test", ex...)
			if err != nil {
				return "", err
			}
			xKey := strings.TrimPrefix(strings.TrimSuffix(util.ToJSONCompact(ex), "]"), "[")
			res[xKey] = x
		}
		ps.Title = "Scripting"
		ps.Data = map[string]any{"script": src, "results": res}
		return controller.Render(rc, as, &vscripting.Detail{Path: key, Script: src, Results: res}, ps, "scripting", key)
	})
}
